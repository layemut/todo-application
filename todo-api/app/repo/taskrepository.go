package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
)

type TaskRepository interface {
	FindAll(project *model.Project) ([]*model.Task, error)
	FindByID(project *model.Project, ID int) (*model.Task, error)
	Create(task *model.Task) error
	Update(task *model.Task) error
	Delete(task *model.Task) error
}

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func (tri *TaskRepositoryImpl) FindAll(project *model.Project) ([]*model.Task, error) {
	if err := tri.DB.Model(project).Related(&project.Tasks).Error; err != nil {
		return nil, err
	}

	return project.Tasks, nil
}

func (tri *TaskRepositoryImpl) FindByID(project *model.Project, ID int) (*model.Task, error) {
	var task *model.Task

	if err := tri.DB.Model(project).Where("ID = ?", ID).Related(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (tri *TaskRepositoryImpl) Create(task *model.Task) error {
	if err := tri.DB.Save(task).Error; err != nil {
		return err
	}
	return nil
}

func (tri *TaskRepositoryImpl) Delete(task *model.Task) error {
	if err := tri.DB.Delete(&task).Error; err != nil {
		return err
	}

	return nil
}

func (tri *TaskRepositoryImpl) Update(task *model.Task) error {
	if err := tri.DB.Save(task).Error; err != nil {
		return err
	}
	return nil
}