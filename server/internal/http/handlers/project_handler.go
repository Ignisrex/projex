package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"project-tracker/server/internal/domain"
	"project-tracker/server/internal/services"
)

type ProjectHandler struct {
	projectService services.ProjectService
	taskService    services.TaskService
}

func NewProjectHandler(projectService services.ProjectService, taskService services.TaskService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		taskService:    taskService,
	}
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	project, err := h.projectService.GetProject(id)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var payload domain.Project
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}
	project, err := h.projectService.CreateProject(payload)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload domain.Project
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}
	project, err := h.projectService.UpdateProject(id, payload)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.projectService.DeleteProject(id); err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) ListTasksByProject(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	tasks, err := h.projectService.ListTasksByProject(projectID)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *ProjectHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "id")
	var payload domain.Task
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}
	task, err := h.projectService.CreateTask(projectID, payload)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *ProjectHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var payload domain.Task
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request payload"})
		return
	}
	task, err := h.taskService.UpdateTask(id, payload)
	if err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *ProjectHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.taskService.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusNotImplemented, map[string]string{"error": "not implemented"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
