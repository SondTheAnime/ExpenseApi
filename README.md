# API de Controle de Despesas

Uma API RESTful para controle de despesas pessoais desenvolvida em Go.

## Tecnologias Utilizadas

- Go 1.23+
- JWT para autenticação
- Gin
- Gorm
- Docker
- Scalar
- PostgreSQL

## Estrutura do Projeto

```
ExpenseApi/
├── cmd/api/           # Ponto de entrada da aplicação
├── internal/          # Código interno da aplicação
│   ├── auth/         # Autenticação e JWT
│   ├── config/       # Configurações
│   ├── handler/      # Handlers HTTP
│   ├── middleware/   # Middlewares
│   ├── model/        # Modelos de dados
│   ├── repository/   # Camada de acesso a dados
│   └── service/      # Lógica de negócios
└── pkg/              # Pacotes reutilizáveis
    └── validator/    # Validação de dados
```

## Funcionalidades

- Autenticação de usuários (registro e login)
- CRUD de despesas
- Filtros por período
- Categorização de despesas
- Proteção de rotas com JWT

## Como Executar

1. Clone o repositório
2. Configure as variáveis de ambiente necessárias
3. Execute `go mod download` para baixar as dependências
4. Execute `go run cmd/api/main.go` para iniciar o servidor

## Endpoints da API

### Autenticação
- POST /api/v1/auth/register - Registro de novo usuário
- POST /api/v1/auth/login - Login de usuário

### Despesas
- GET /api/v1/expenses - Lista todas as despesas
- POST /api/v1/expenses - Cria uma nova despesa
- GET /api/v1/expenses/{id} - Obtém uma despesa específica
- PUT /api/v1/expenses/{id} - Atualiza uma despesa
- DELETE /api/v1/expenses/{id} - Remove uma despesa

## Categorias de Despesas

- Mantimentos
- Lazer
- Eletrônica
- Utilitários
- Roupas
- Saúde
- Outros 