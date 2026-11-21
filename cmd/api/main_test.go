package main

import (
	"sync"
	"testing"
)

func TestRepositorySnapshots(t *testing.T) {
	repo := NewTaskRepository()
	created := repo.Create("Revisar")
	list := repo.List()
	list[0].Name = "alterado"
	if repo.List()[0].Name != "Revisar" {
		t.Fatal("list exposes repository storage")
	}
	processing := repo.FetchPendingAndMarkProcessing()
	if processing == nil || processing.Status != "processing" {
		t.Fatal("task was not claimed")
	}
	repo.MarkDone(created.ID)
	if processing.Status != "processing" || repo.List()[0].Status != "done" {
		t.Fatal("worker snapshot changed after repository update")
	}
	if repo.FetchPendingAndMarkProcessing() != nil {
		t.Fatal("completed task claimed again")
	}
}

func TestConcurrentRepositoryAccess(t *testing.T) {
	repo := NewTaskRepository()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				repo.Create("tarefa")
				for _, task := range repo.List() {
					_ = task.Status
				}
				if task := repo.FetchPendingAndMarkProcessing(); task != nil {
					repo.MarkDone(task.ID)
				}
			}
		}()
	}
	wg.Wait()
	if len(repo.List()) != 200 {
		t.Fatal("tasks lost")
	}
}
