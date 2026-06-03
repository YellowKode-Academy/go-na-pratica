# Stage 1: builder — compilar o binário
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copiar apenas os arquivos de módulos primeiro (cache de dependências)
COPY go.mod go.sum ./
RUN go mod download

# Copiar o código fonte
COPY . .

# Compilar o binário estático
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
      -ldflags "-s -w" \
      -o /linkvault-api \
      ./cmd/api

# Criar diretório de dados com ownership do nonroot (uid=65532)
RUN mkdir -p /data && chown 65532:65532 /data


# Stage 2: runner — imagem final mínima
FROM gcr.io/distroless/static-debian12

# Copiar o binário e o diretório de dados do stage de build
COPY --from=builder /linkvault-api /linkvault-api
COPY --from=builder /data /data

# Usuário não-root por segurança
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/linkvault-api"]
