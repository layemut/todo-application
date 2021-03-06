package app

import (
	"github.com/layemut/todo-application/todo-api/app/controller"
	"github.com/layemut/todo-application/todo-api/app/model"
	"github.com/layemut/todo-application/todo-api/app/repo"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/layemut/todo-application/todo-api/config"
)

// App has router and db instances
type App struct {
	Router            *mux.Router
	DB                *gorm.DB
	ProjectController controller.ProjectController
	TaskController    controller.TaskController
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	db, err := gorm.Open("sqlite3", "todo.db")
	if err != nil {
		log.Fatal("failed to connect database")
	}

	a.DB = model.DBMigrate(db)
	a.ProjectController = controller.ProjectController{ProjectRepository: &repo.ProjectRepositoryImpl{DB: db}}
	a.TaskController = controller.TaskController{ProjectRepository: &repo.ProjectRepositoryImpl{DB: db}, TaskRepository: &repo.TaskRepositoryImpl{DB: db}}
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/projects", a.ProjectController.GetAllProjects)
	a.Post("/projects", a.ProjectController.CreateProject)
	a.Get("/projects/{title}", a.ProjectController.GetProject)
	a.Put("/projects/{title}", a.ProjectController.UpdateProject)
	a.Delete("/projects/{title}", a.ProjectController.DeleteProject)

	// Routing for handling the tasks
	a.Get("/projects/{title}/tasks", a.TaskController.GetAllTasks)
	a.Post("/projects/{title}/tasks", a.TaskController.CreateTask)
	a.Get("/projects/{title}/tasks/{id:[0-9]+}", a.TaskController.GetTask)
	a.Put("/projects/{title}/tasks/{id:[0-9]+}", a.TaskController.UpdateTask)
	a.Delete("/projects/{title}/tasks/{id:[0-9]+}", a.TaskController.DeleteTask)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
