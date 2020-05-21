package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
)

type ProjectRepository interface {
	FindAll() ([]*model.Project, error)
	FindByTitle(title string) (*model.Project, error)
	Create(project *model.Project) error
	Update(project *model.Project) error
	Delete(project *model.Project) error
}

type ProjectRepositoryImpl struct {
	DB *gorm.DB
}

func (pri *ProjectRepositoryImpl) FindAll() ([]*model.Project, error) {
	var projects []*model.Project

	if err := pri.DB.Find(&projects).Error; err != nil {
		return nil, err
	}

	for _, p := range projects {
		if err := pri.DB.Model(&p).Related(&p.Tasks).Error; err != nil {
			return nil, err
		}
	}

	return projects, nil
}

func (pri *ProjectRepositoryImpl) FindByTitle(title string) (*model.Project, error) {
	var project *model.Project

	if err := pri.DB.Find(&project, model.Project{Title: title}).Related(&project.Tasks).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (pri *ProjectRepositoryImpl) Create(project *model.Project) error {
	if err := pri.DB.Save(&project).Error; err != nil {
		return err
	}

	return nil
}

func (pri *ProjectRepositoryImpl) Update(project *model.Project) error {
	return pri.Create(project)
}

func (pri *ProjectRepositoryImpl) Delete(project *model.Project) error {
	if err := pri.DB.Delete(&project).Error; err != nil {
		return err
	}

	return nil
}

