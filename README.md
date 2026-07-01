# LinkVault

[![Capa do Livro](cover.jpg)](https://www.amazon.com.br/dp/B0H3WTFPCW)

> **[Disponivel na Amazon.com.br](https://www.amazon.com.br/dp/B0H3WTFPCW)** — R$24,97

---

Repositorio oficial do livro **Go na Prática**.

Gerenciador de links pessoal com API REST e CLI.

## Instalacao

```bash
go install github.com/yellowkode-academy/linkvault/cmd/cli@latest
```

## Uso

```bash
# API
go run ./cmd/api

# CLI
go run ./cmd/cli add https://go.dev --title "Go oficial" --tags go,docs
go run ./cmd/cli list
go run ./cmd/cli search go
```

## Variaveis de ambiente

| Variavel       | Padrao         | Descricao                          |
|----------------|----------------|------------------------------------|
| `DATABASE_URL` | `linkvault.db` | Caminho do banco SQLite (API)      |
| `LINKVAULT_DB` | `linkvault.db` | Caminho do banco SQLite (CLI)      |
| `PORT`         | `8080`         | Porta da API REST                  |

## Endpoints

| Metodo   | Rota           | Descricao                     |
|----------|----------------|-------------------------------|
| GET      | /health        | Health check                  |
| GET      | /links         | Lista todos os links          |
| GET      | /links?q=query | Busca por texto               |
| POST     | /links         | Cria um novo link             |
| GET      | /links/{id}    | Retorna link por ID           |
| DELETE   | /links/{id}    | Remove um link                |

## Testes

```bash
go test ./...
```

## Estrutura

```
linkvault/
├── cmd/
│   ├── api/          # servidor HTTP
│   └── cli/          # linha de comando
├── internal/
│   ├── link/         # dominio: struct Link e validacao
│   ├── storage/      # repositorios: memoria e SQLite
│   ├── api/          # handlers HTTP
│   └── middleware/   # logger e CORS
└── go.mod
```
