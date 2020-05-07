package model

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Project struct {
	gorm.Model
	Title    string `gorm:"unique" json:"title"`
	Archived bool   `json:"archived"`
	Tasks    []Task `gorm:"foreignkey:ID" json:"tasks"`
}

func (p *Project) Archive() {
	p.Archived = true
}

func (p *Project) Restore() {
	p.Archived = false
}

type Task struct {
	gorm.Model
	Title     string     `json:"title"`
	Priority  string     `json:"priority"`
	Deadline  *time.Time `gorm:"default:null" json:"deadline"`
	Done      bool       `json:"done"`
	ProjectID uint       `json:"project_id"`
}

func (t *Task) Complete() {
	t.Done = true
}

func (t *Task) Undo() {
	t.Done = false
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.CreateTable(&Project{})
	db.CreateTable(&Task{})
	db.AutoMigrate(&Task{}, &Project{})
	return db
}
