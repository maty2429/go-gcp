package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"task-api/internal/repository"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

type createTaskRequest struct {
	Title string `json:"title"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.index)
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("GET /tasks/{id}", h.getTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.deleteTask)
}

func (h *TaskHandler) index(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "task-api is running",
		"status":  "ok",
	})
}

func (h *TaskHandler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	if h.repo == nil {
		writeDatabaseUnavailable(w)
		return
	}

	tasks, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not list tasks")
		return
	}

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	if h.repo == nil {
		writeDatabaseUnavailable(w)
		return
	}

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	task, err := h.repo.Create(r.Context(), title)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create task")
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	if h.repo == nil {
		writeDatabaseUnavailable(w)
		return
	}

	id, ok := parseID(w, r)
	if !ok {
		return
	}

	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "could not get task")
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	if h.repo == nil {
		writeDatabaseUnavailable(w)
		return
	}

	id, ok := parseID(w, r)
	if !ok {
		return
	}

	err := h.repo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "could not delete task")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return "", false
	}

	return id, true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeDatabaseUnavailable(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "database not configured yet")
}
