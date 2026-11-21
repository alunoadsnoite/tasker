# Tasker

**API de tarefas com workers concorrentes e cliente de terminal.**

`Go · HTTP · Goroutines`

[Reportar problema](https://github.com/alunoadsnoite/tasker/issues)

Projeto de estudo de backend em Go: crie tarefas por HTTP ou pelo menu do terminal e acompanhe o processamento por um conjunto de workers.

## Recursos

- API HTTP para criar e listar tarefas.
- Filtro por status: `pending`, `processing` e `done`.
- Número de workers configurável na inicialização.
- CLI interativa para criar, listar e acompanhar tarefas e workers.
- Armazenamento em memória: os dados são perdidos ao reiniciar a API.

## Executar

Instale a versão de Go indicada em [go.mod](go.mod) e rode na raiz do repositório:

```sh
go mod download
go run ./cmd/api -workers 3
```

Em outro terminal:

```sh
go run ./cmd/cli
```

A API escuta na porta 8080. O cliente usa `http://localhost:8080`.

## Experimentar a API

```sh
curl -X POST http://localhost:8080/tasks \
  -H 'Content-Type: application/json' \
  -d '{"name":"Revisar documentação"}'

curl http://localhost:8080/tasks
curl 'http://localhost:8080/tasks?status=done'
curl http://localhost:8080/workers
```

| Método | Endpoint | Função |
|---|---|---|
| POST | `/tasks` | Criar uma tarefa com o campo `name` |
| GET | `/tasks` | Listar tarefas; aceita `?status=pending`, `processing` ou `done` |
| GET | `/workers` | Consultar o estado dos workers |

## Organização

- `cmd/api/`: API e inicialização dos workers.
- `cmd/cli/`: interface interativa de terminal.
- `internal/`: modelos e implementações de apoio.

As entradas da aplicação ficam separadas em `cmd/api/` e `cmd/cli/`. Os arquivos antigos da raiz foram removidos para permitir `go build ./...` sem conflito entre funções `main()`.

## Limites e evolução

Este é um projeto de aprendizado. Antes de disponibilizar a API em produção, é necessário revisar validação de entrada, acesso concorrente aos estados e às listas, persistência, timeouts e tratamento de erros. Os testes automatizados verificam snapshots do repositório e acesso concorrente, com detector de condições de corrida.

## Como contribuir

Abra uma issue com uma proposta ou um problema reproduzível. Em pull requests, explique a mudança e como foi validada.

## Verificações automáticas

O GitHub Actions executa `go vet ./...`, `go test -race ./...` e `go build ./...` nos pushes e pull requests.
