# Estágio de build
FROM golang:1.23-alpine3.21 AS builder

WORKDIR /app

# Copiando os arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiando o código fonte
COPY . .

# Compilando a aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/api ./cmd/api

# Estágio final
FROM alpine:3.21

WORKDIR /app

# Copiando o binário compilado
COPY --from=builder /app/api .
# Copiando a pasta de documentação
COPY --from=builder /app/docs ./docs

# Expondo a porta da API
EXPOSE 8081

# Executando a API
CMD ["./api"] 