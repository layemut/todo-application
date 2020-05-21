package controller

import (
	"github.com/layemut/todo-application/todo-api/app/repo"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
	"github.com/layemut/todo-application/todo-api/app/util"
)

type TaskController struct {
	ProjectRepository repo.ProjectRepository
	TaskRepository    repo.TaskRepository
}

func (tc TaskController) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project, err := tc.ProjectRepository.FindByTitle(title)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	tasks, err := tc.TaskRepository.FindAll(project)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	tasksResponse.Tasks = tasks
	tasksResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, tasksResponse)
}

func (tc TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project, err := tc.ProjectRepository.FindByTitle(title)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	task := model.Task{ProjectID: project.ID}

	if err := task.Parse(w, r); err != nil {
		tasksResponse.Response = model.BadRequestResponse(err)
		RespondJSON(w, http.StatusBadRequest, tasksResponse)
		return
	}

	tasksResponse.Tasks = []*model.Task{&task}
	tasksResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusCreated, tasksResponse)
}

func (tc TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project, err := tc.ProjectRepository.FindByTitle(title)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))

	task, err := tc.TaskRepository.FindByID(project, id)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "Task not found", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	tasksResponse.Tasks = []*model.Task{task}
	tasksResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, tasksResponse)
}

func (tc TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project, err := tc.ProjectRepository.FindByTitle(title)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))

	task, err := tc.TaskRepository.FindByID(project, id)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "task with id: "+strconv.Itoa(id)+" not found", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	if err := task.Parse(w, r); err != nil {
		tasksResponse.Response = model.BadRequestResponse(err)
		RespondJSON(w, http.StatusBadRequest, tasksResponse)
		return
	}

	if err := tc.TaskRepository.Create(task); err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error updating task", err.Error())
		RespondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}

	RespondJSON(w, http.StatusOK, task)
}

func (tc TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project, err := tc.ProjectRepository.FindByTitle(title)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))

	task, err := tc.TaskRepository.FindByID(project, id)
	if err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "task with id: "+strconv.Itoa(id)+" not found", err.Error())
		RespondJSON(w, http.StatusOK, tasksResponse)
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error deleting task", err.Error())
		RespondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}
	tasksResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, tasksResponse)
}

// getTaskOr404 gets a task instance if exists, or respond the 404 error otherwise
func (tc TaskController) getTaskOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Task {
	tasksResponse := model.TasksResponse{}
	task := model.Task{}
	if err := db.First(&task, id).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "Task not found", err.Error())
		RespondJSON(w, http.StatusNotFound, tasksResponse)
		return nil
	}
	return &task
}
