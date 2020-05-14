package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
)

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}
	projects := []*model.Project{}
	db.Find(&projects)

	projectsResponse.Projects = projects

	if len(projects) == 0 {
		projectsResponse.Response = model.PrepareResponse(400, "No Project Found.", nil)
		respondJSON(w, http.StatusOK, projectsResponse)
		return
	}
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, projectsResponse)
}

func CreateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}
	project := model.Project{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondJSON(w, http.StatusInternalServerError, model.PrepareResponse(500, "Error saving project", err))
		return
	}

	projects := []*model.Project{}
	projects = append(projects, &project)
	projectsResponse.Projects = projects
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusCreated, project)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}
	vars := mux.Vars(r)

	title := vars["title"]
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}
	projects := []*model.Project{project}
	projectsResponse.Projects = projects
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, projectsResponse)
}

func UpdateProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	title := vars["title"]
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	title := vars["title"]
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}
	if err := db.Delete(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	title := vars["title"]
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}
	project.Archive()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	title := vars["title"]
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}
	project.Restore()
	if err := db.Save(&project).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, project)
}

// getProjectOr404 gets a project instance if exists, or respond the 404 error otherwise
func getProjectOr404(db *gorm.DB, title string, w http.ResponseWriter, r *http.Request) *model.Project {
	projectsResponse := model.ProjectsResponse{}
	project := model.Project{}
	if err := db.First(&project, model.Project{Title: title}).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(404, title+" not found.", err)
		respondJSON(w, http.StatusOK, projectsResponse)
		return nil
	}
	return &project
}
