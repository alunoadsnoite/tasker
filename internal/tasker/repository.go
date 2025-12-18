package tasker

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskRepository struct {
	mu    sync.Mutex
	tasks []Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: []Task{},
	}
}

func (r *TaskRepository) Create(name string) Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	task := Task{
		ID:        uuid.NewString(),
		Name:      name,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	r.tasks = append(r.tasks, task)
	return task
}

func (r *TaskRepository) List() []Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.tasks
}

// retorna a próxima task pendente e marca como processing
func (r *TaskRepository) FetchPendingAndMarkProcessing() *Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.tasks {
		if r.tasks[i].Status == "pending" {
			r.tasks[i].Status = "processing"
			return &r.tasks[i]
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
