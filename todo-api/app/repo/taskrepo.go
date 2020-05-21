package repo

import (
	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
	"net/http"
	"strconv"
)

type TaskRepository interface {
	FindAll(project *model.Project) ([]*model.Task, error)
	FindByID(project *model.Project, ID int) (*model.Task, error)
	Create(ID int) error
	Update(ID int) error
	Delete(ID int) error
}

type TaskRepositoryImpl struct {
	DB *gorm.DB
}

func (tri *TaskRepositoryImpl) FindAll(project *model.Project) ([]*model.Task, error) {
	var tasks []*model.Task

	if err := tri.DB.Model(project).Related(tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tri *TaskRepositoryImpl) FindByID(project *model.Project, ID int) (*model.Task, error) {
	var task *model.Task

	if err := tri.DB.Model(project).Where("ID = ?", ID).Related(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (tri *TaskRepositoryImpl) Create(task *model.Task) error {
	task, err := tc.TaskRepository.FindByID(project, id)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "task with id: "+strconv.Itoa(id)+" not found", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}
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
