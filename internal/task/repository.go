package tasker

import "sync"

type TaskRepository struct {
	mu    sync.Mutex
	tasks []Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make([]Task, 0),
	}
}

func (r *TaskRepository) Add(task Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, task)
}

func (r *TaskRepository) List() []Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]Task{}, r.tasks...)
}

func (r *TaskRepository) GetByID(id string) (Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range r.tasks {
		if t.ID == id {
			return t, true
		}
	}
	return Task{}, false
}
