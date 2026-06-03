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


# Stage 2: runner — imagem final mínima
FROM alpine:3.21

# Criar usuário não-root e diretório de dados
RUN addgroup -g 65532 nonroot && \
    adduser -u 65532 -G nonroot -D nonroot && \
    mkdir -p /data && chown nonroot:nonroot /data

# Copiar apenas o binário do stage de build
COPY --from=builder /linkvault-api /linkvault-api

# Usuário não-root por segurança
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/linkvault-api"]
