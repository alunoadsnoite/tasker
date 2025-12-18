package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiURL = "http://localhost:8080"

type createTaskRequest struct {
	Name string `json:"name"`
}

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "create":
		createCommand(os.Args[2:])
	case "list":
		listCommand()
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("uso:")
	fmt.Println("  go run cli.go create \"nome da tarefa\"")
	fmt.Println("  go run cli.go list")
}

func createCommand(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("erro: informe o nome da tarefa")
		return
	}

	taskName := fs.Arg(0)

	req := createTaskRequest{Name: taskName}
	data, _ := json.Marshal(req)

	resp, err := http.Post(
		apiURL+"/tasks",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		fmt.Println("erro ao chamar API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("erro:", resp.Status)
		return
	}

	fmt.Println("Task criada:", taskName)
}

func listCommand() {
	resp, err := http.Get(apiURL + "/tasks")
	if err != nil {
		fmt.Println("erro ao chamar API:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var tasks []Task
	json.Unmarshal(body, &tasks)

	if len(tasks) == 0 {
		fmt.Println("nenhuma task encontrada")
		return
	}

	for _, t := range tasks {
		fmt.Printf("- [%s] %s (%s)\n", t.Status, t.Name, t.ID)
	}
}
