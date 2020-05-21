package controller

import (
	"github.com/layemut/todo-application/todo-api/app/repo"
	"net/http"

	"github.com/layemut/todo-application/todo-api/app/model"
	"github.com/layemut/todo-application/todo-api/app/util"
)

type ProjectController struct {
	ProjectRepository repo.ProjectRepository
}

func (pc ProjectController) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	projects, err := pc.ProjectRepository.FindAll()

	if len(projects) == 0 || err != nil {
		projectsResponse.Response = model.PrepareResponse(400, "No Project Found.", err.Error())
		RespondJSON(w, http.StatusOK, projectsResponse)
		return
	}

	projectsResponse.Projects = projects
	projectsResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, projectsResponse)
}

func (pc ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}
	project := model.Project{}

	if err := project.Parse(w, r); err != nil {
		projectsResponse.Response = model.BadRequestResponse(err)
		RespondJSON(w, http.StatusBadRequest, projectsResponse)
		return
	}

	if err := pc.ProjectRepository.Create(&project); err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error saving project", err.Error())
		RespondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	var projects []*model.Project
	projects = append(projects, &project)
	projectsResponse.Projects = projects
	projectsResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusCreated, projectsResponse)
}

func (pc ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")

	project, err := pc.ProjectRepository.FindByTitle(title)
	if err != nil {
		projectsResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, projectsResponse)
		return
	}

	projectsResponse.Projects = []*model.Project{project}
	projectsResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, projectsResponse)
}

func (pc ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")

	project, err := pc.ProjectRepository.FindByTitle(title)
	if err != nil {
		projectsResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, projectsResponse)
		return
	}

	if err := project.Parse(w, r); err != nil {
		projectsResponse.Response = model.BadRequestResponse(err)
		RespondJSON(w, http.StatusBadRequest, projectsResponse)
		return
	}

	if err := pc.ProjectRepository.Update(project); err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error updating project", err.Error())
		RespondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projectsResponse.Projects = []*model.Project{project}
	projectsResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusOK, projectsResponse)
}

func (pc ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	projectsResponse := model.ProjectsResponse{}

	title := util.GetParam(r, "title")

	project, err := pc.ProjectRepository.FindByTitle(title)
	if err != nil {
		projectsResponse.Response = model.PrepareResponse(404, "project with title: "+title+" not found.", err.Error())
		RespondJSON(w, http.StatusOK, projectsResponse)
		return
	}

	if err := pc.ProjectRepository.Delete(project); err != nil {
		projectsResponse.Response = model.PrepareResponse(500, "Error deleting project", err.Error())
		RespondJSON(w, http.StatusInternalServerError, projectsResponse)
		return
	}

	projectsResponse.Response = model.SuccessResponse()
	RespondJSON(w, http.StatusNoContent, projectsResponse)
}
