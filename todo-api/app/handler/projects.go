package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/layemut/todo-application/todo-api/app/model"
	"github.com/layemut/todo-application/todo-api/app/util"
)

func GetAllProjects(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}
	projects := []*model.Project{}
	db.Find(&projects)

	for _, p := range projects {
		db.Model(&p).Related(&p.Tasks)
	}
	
	projectsResponse.Projects = projects

	if len(projects) == 0 {
		projectsResponse.Response = model.PrepareResponse(400, "No Project Found.", "")
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
		projectsResponse.Response = model.PrepareResponse(500, "Bad Request", err.Error())
		respondJSON(w, http.StatusBadRequest, projectsResponse)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error saving project", err.Error())
		respondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projects := []*model.Project{}
	projects = append(projects, &project)
	projectsResponse.Projects = projects
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusCreated, projectsResponse)
}

func GetProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")
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
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&project); err != nil {
		projectsResponse.Response = model.PrepareResponse(400, "Bad Request", err.Error())
		respondJSON(w, http.StatusBadRequest, projectsResponse)
		return
	}
	defer r.Body.Close()

	if err := db.Save(&project).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error updating project", err.Error())
		respondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, projectsResponse)
}

func DeleteProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	if err := db.Delete(&project).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error deleting project", err.Error())
		respondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusNoContent, projectsResponse)
}

func ArchiveProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	project.Archive()

	if err := db.Save(&project).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error archiving project", err.Error())
		respondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projectsResponse.Projects = []*model.Project{project}
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, projectsResponse)
}

func RestoreProject(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")
	project := getProjectOr404(db, title, w, r)
	if project == nil {
		return
	}

	project.Restore()

	if err := db.Save(&project).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error restoring project", err.Error())
		respondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projectsResponse.Projects = []*model.Project{project}
	projectsResponse.Response = model.SuccessResponse()
	respondJSON(w, http.StatusOK, projectsResponse)
}

// getProjectOr404 gets a project instance if exists, or respond the 404 error otherwise
func getProjectOr404(db *gorm.DB, title string, w http.ResponseWriter, r *http.Request) *model.Project {
	projectsResponse := model.ProjectsResponse{}
	project := model.Project{}
	tasks := []model.Task{}
	if err := db.First(&project, model.Project{Title: title}).Related(&tasks).Error; err != nil {
		projectsResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		respondJSON(w, http.StatusOK, projectsResponse)
		return nil
	}
	project.Tasks = tasks
	return &project
}
