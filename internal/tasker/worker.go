package tasker

import (
	"log"
	"time"
)

// StartWorkers inicia N workers concorrentes para processar tasks
func StartWorkers(repo *TaskRepository, workerCount int) {
	taskChan := make(chan *Task, 100) // canal buffer

	// Inicializa N workers
	for i := 0; i < workerCount; i++ {
		go func(id int) {
			for task := range taskChan {
				log.Printf("[Worker %d] Processando task: %s\n", id, task.Name)
				time.Sleep(3 * time.Second) // simula processamento
				repo.MarkDone(task.ID)
				log.Printf("[Worker %d] Task concluída: %s\n", id, task.Name)
			}
		}(i + 1)
	}

	// Goroutine que empurra tasks pendentes para o canal
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
