package http

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"time"

	task "tasker/internal/task"
)

type TaskHandler struct {
	Repo *task.TaskRepository
}

type createTaskRequest struct {
	Name string `json:"name"`
}

func NewTaskHandler(repo *task.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

func (h *TaskHandler) HandleTasks(w nethttp.ResponseWriter, r *nethttp.Request) {
	switch r.Method {

	case nethttp.MethodPost:
		var req createTaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			nethttp.Error(w, "JSON inválido", nethttp.StatusBadRequest)
			return
		}

		if req.Name == "" {
			nethttp.Error(w, "name é obrigatório", nethttp.StatusBadRequest)
			return
		}

		t := task.Task{
			ID:        fmt.Sprint(time.Now().UnixNano()),
			Name:      req.Name,
			Status:    "pending",
			CreatedAt: time.Now(),
		}

		h.Repo.Add(t)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusCreated)
		json.NewEncoder(w).Encode(t)

	case nethttp.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.Repo.List())

	default:
		w.WriteHeader(nethttp.StatusMethodNotAllowed)
	}
}
