package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"tasker/internal/tasker"
	"time"

	"github.com/google/uuid"
)

type TaskRepository struct {
	mu    sync.Mutex
	tasks []tasker.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: []tasker.Task{}}
}

func (r *TaskRepository) Create(name string) tasker.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	t := tasker.Task{
		ID:        uuid.NewString(),
		Name:      name,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	r.tasks = append(r.tasks, t)
	return t
}

func (r *TaskRepository) List() []tasker.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]tasker.Task{}, r.tasks...)
}

func (r *TaskRepository) FetchPendingAndMarkProcessing() *tasker.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.tasks {
		if r.tasks[i].Status == "pending" {
			r.tasks[i].Status = "processing"
			task := r.tasks[i]
			return &task
		}
	}
	return nil
}

func (r *TaskRepository) MarkDone(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.tasks {
		if r.tasks[i].ID == id {
			r.tasks[i].Status = "done"
			return
		}
	}
}

type WorkerStatus struct {
	ID     int    `json:"id"`
	TaskID string `json:"task_id,omitempty"`
	Name   string `json:"task_name,omitempty"`
	Busy   bool   `json:"busy"`
}

func StartWorkers(repo *TaskRepository, n int, statusChan chan WorkerStatus) {
	taskChan := make(chan *tasker.Task, 100)

	for i := 0; i < n; i++ {
		go func(id int) {
			for task := range taskChan {
				statusChan <- WorkerStatus{ID: id, TaskID: task.ID, Name: task.Name, Busy: true}
				log.Printf("[Worker %d] Processando task: %s\n", id, task.Name)
				time.Sleep(3 * time.Second)
				repo.MarkDone(task.ID)
				statusChan <- WorkerStatus{ID: id, TaskID: task.ID, Name: task.Name, Busy: false}
				log.Printf("[Worker %d] Task concluída: %s\n", id, task.Name)
			}
		}(i + 1)
	}

	go func() {
		for {
			task := repo.FetchPendingAndMarkProcessing()
			if task != nil {
				taskChan <- task
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func main() {
	workers := flag.Int("workers", 3, "Número de workers concorrentes")
	flag.Parse()

	repo := NewTaskRepository()
	statusChan := make(chan WorkerStatus, *workers*10)
	StartWorkers(repo, *workers, statusChan)

	workerStatuses := make(map[int]WorkerStatus)
	var statusesMu sync.RWMutex
	go func() {
		for s := range statusChan {
			statusesMu.Lock()
			workerStatuses[s.ID] = s
			statusesMu.Unlock()
		}
	}()

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req struct {
				Name string `json:"name"`
			}
			json.NewDecoder(r.Body).Decode(&req)
			task := repo.Create(req.Name)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(task)

		case http.MethodGet:
			tasks := repo.List()
			statusFilter := r.URL.Query().Get("status")
			if statusFilter != "" {
				filtered := []tasker.Task{}
				for _, t := range tasks {
					if t.Status == statusFilter {
						filtered = append(filtered, t)
					}
				}
				tasks = filtered
			}
			json.NewEncoder(w).Encode(tasks)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/workers", func(w http.ResponseWriter, r *http.Request) {
		statuses := []WorkerStatus{}
		statusesMu.RLock()
		for _, s := range workerStatuses {
			statuses = append(statuses, s)
		}
		statusesMu.RUnlock()
		json.NewEncoder(w).Encode(statuses)
	})

	log.Printf("API rodando em http://localhost:8080 com %d workers\n", *workers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
