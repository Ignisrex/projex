package http

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"project-tracker/server/internal/config"
	"project-tracker/server/internal/http/handlers"
	"project-tracker/server/internal/services"
)

func NewRouter(cfg config.Config, projectService services.ProjectService, taskService services.TaskService) stdhttp.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	healthHandler := handlers.NewHealthHandler()
	projectHandler := handlers.NewProjectHandler(projectService, taskService)

	r.Get("/health", healthHandler.GetHealth)
	r.Get("/projects", projectHandler.ListProjects)
	r.Get("/projects/{id}", projectHandler.GetProject)
	r.Post("/projects", projectHandler.CreateProject)
	r.Put("/projects/{id}", projectHandler.UpdateProject)
	r.Delete("/projects/{id}", projectHandler.DeleteProject)
	r.Get("/projects/{id}/tasks", projectHandler.ListTasksByProject)
	r.Post("/projects/{id}/tasks", projectHandler.CreateTask)
	r.Put("/tasks/{id}", projectHandler.UpdateTask)
	r.Delete("/tasks/{id}", projectHandler.DeleteTask)

	return r
}
