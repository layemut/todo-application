package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
	"github.com/layemut/todo-application/todo-api/app/util"
)

func GetAllTasks(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	tasks := []*model.Task{}
	if err := db.Model(&project).Related(&tasks).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error getting tasks", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}

	tasksResponse.Response = model.SuccessResponse()
	tasksResponse.Tasks = tasks
	respondJSON(w, http.StatusOK, tasksResponse)
}

func CreateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")

	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	task := model.Task{ProjectID: project.ID}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		tasksResponse.Response = model.BadRequestResponse(err)
		respondJSON(w, http.StatusBadRequest, tasksResponse)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&task).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error creating taks", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}

	tasks := []*model.Task{}
	tasks = append(tasks, &task)
	tasksResponse.Response = model.SuccessResponse()
	tasksResponse.Tasks = tasks
	respondJSON(w, http.StatusCreated, tasksResponse)
}

func GetTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	tasksResponse.Response = model.SuccessResponse()
	tasksResponse.Tasks = []*model.Task{task}
	respondJSON(w, http.StatusOK, tasksResponse)
}

func UpdateTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		tasksResponse.Response = model.BadRequestResponse(err)
		respondJSON(w, http.StatusBadRequest, tasksResponse)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&task).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error updating task", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}
	respondJSON(w, http.StatusOK, task)
}

func DeleteTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error deleting task", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}
	tasksResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, tasksResponse)
}

func CompleteTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	task.Complete()
	if err := db.Save(&task).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error completing task", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}
	tasksResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, tasksResponse)
}

func UndoTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	tasksResponse := model.TasksResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	id, _ := strconv.Atoi(util.GetParam(r, "id"))
	task := getTaskOr404(db, id, w, r)
	if task == nil {
		return
	}

	task.Undo()
	if err := db.Save(&task).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(500, "Error undoing task", err.Error())
		respondJSON(w, http.StatusInternalServerError, tasksResponse)
		return
	}

	tasksResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, tasksResponse)
}

// getTaskOr404 gets a task instance if exists, or respond the 404 error otherwise
func getTaskOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.Task {
	tasksResponse := model.TasksResponse{}
	task := model.Task{}
	if err := db.First(&task, id).Error; err != nil {
		tasksResponse.Response = model.PrepareResponse(404, "Task not found", err.Error())
		respondJSON(w, http.StatusNotFound, tasksResponse)
		return nil
	}
	return &task
}
