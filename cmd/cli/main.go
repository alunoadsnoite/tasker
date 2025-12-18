package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"tasker/internal/tasker"
	"time"
)

const apiTasksURL = "http://localhost:8080/tasks"
const apiWorkersURL = "http://localhost:8080/workers"

func createTask(name string) {
	payload := map[string]string{"name": name}
	data, _ := json.Marshal(payload)
	resp, err := http.Post(apiTasksURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Erro ao criar task:", err)
		return
	}
	defer resp.Body.Close()
	var t tasker.Task
	json.NewDecoder(resp.Body).Decode(&t)
	fmt.Printf("Task criada: %s (ID: %s)\n", t.Name, t.ID)
}

func listTasks(filter string) {
	url := apiTasksURL
	if filter != "" {
		url += "?status=" + filter
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Erro ao listar tasks:", err)
		return
	}
	defer resp.Body.Close()
	var tasks []tasker.Task
	json.NewDecoder(resp.Body).Decode(&tasks)
	if len(tasks) == 0 {
		fmt.Println("Nenhuma task encontrada.")
		return
	}
	for _, t := range tasks {
		fmt.Printf("- %s | %s | %s\n", t.ID, t.Name, t.Status)
	}
}

func listWorkers() {
	resp, err := http.Get(apiWorkersURL)
	if err != nil {
		fmt.Println("Erro ao obter status dos workers:", err)
		return
	}
	defer resp.Body.Close()
	var workers []struct {
		ID     int    `json:"id"`
		TaskID string `json:"task_id,omitempty"`
		Name   string `json:"task_name,omitempty"`
		Busy   bool   `json:"busy"`
	}
	json.NewDecoder(resp.Body).Decode(&workers)
	fmt.Println("=== Workers ===")
	for _, w := range workers {
		status := "livre"
		if w.Busy {
			status = fmt.Sprintf("ocupado com %s", w.Name)
		}
		fmt.Printf("- Worker %d: %s\n", w.ID, status)
	}
}

func watchTasksAndWorkers(filter string) {
	fmt.Println("Acompanhando tasks e workers em tempo real. Ctrl+C para sair.")
	for {
		listTasks(filter)
		listWorkers()
		fmt.Println()
		time.Sleep(2 * time.Second)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== CLI Tasker ===")
		fmt.Println("1 - Criar task")
		fmt.Println("2 - Listar tasks")
		fmt.Println("3 - Acompanhar tasks e workers em tempo real")
		fmt.Println("4 - Listar tasks por status")
		fmt.Println("5 - Acompanhar tasks e workers por status")
		fmt.Println("0 - Sair")
		fmt.Print("Escolha: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("Nome da task: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			createTask(name)

		case "2":
			listTasks("")

		case "3":
			watchTasksAndWorkers("")

		case "4":
			fmt.Print("Status (pending/processing/done): ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)
			listTasks(status)

		case "5":
			fmt.Print("Status (pending/processing/done): ")
			status, _ := reader.ReadString('\n')
			status = strings.TrimSpace(status)
			watchTasksAndWorkers(status)

		case "0":
			fmt.Println("Saindo...")
			return

		default:
			fmt.Println("Opção inválida")
		}

		fmt.Println()
	}
}
