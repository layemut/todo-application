package model

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Project struct {
	gorm.Model
	Title    string `gorm:"unique" json:"title"`
	Archived bool   `json:"archived"`
	Tasks    []Task `gorm:"foreignkey:ProjectID" json:"tasks"`
}

func (p *Project) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

type Task struct {
	gorm.Model
	Title     string     `json:"title"`
	Priority  string     `json:"priority"`
	Deadline  *time.Time `gorm:"default:null" json:"deadline"`
	Done      bool       `json:"done"`
	ProjectID uint       `json:"project_id"`
}

func (t *Task) Parse(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&t); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.CreateTable(&Project{})
	db.CreateTable(&Task{})
	db.AutoMigrate(&Task{}, &Project{})
	return db
}
