# API de Controle de Despesas

Uma API RESTful robusta para controle de despesas pessoais desenvolvida em Go, seguindo as melhores práticas de desenvolvimento e arquitetura limpa.

## 🚀 Tecnologias Utilizadas

- Go 1.23+
- PostgreSQL (Banco de dados principal)
- JWT (JSON Web Tokens para autenticação)
- ServeMux (Router nativo do Go 1.22+)
- Docker & Docker Compose (Containerização)
- Make (Automação de comandos)
- Swagger/OpenAPI (Documentação da API)

## 📁 Estrutura do Projeto

```
ExpenseApi/
├── cmd/
│   └── api/                 # Ponto de entrada da aplicação
├── internal/                # Código interno da aplicação
│   ├── auth/               # Autenticação e JWT
│   ├── config/             # Configurações da aplicação
│   ├── handler/            # Handlers HTTP
│   ├── middleware/         # Middlewares (CORS, Auth, etc)
│   ├── model/              # Modelos de dados
│   ├── repository/         # Camada de acesso a dados
│   ├── server/             # Configuração do servidor HTTP
│   └── service/            # Lógica de negócios
├── docs/                    # Documentação OpenAPI/Swagger
│   ├── components/
│   ├── paths/
│   ├── responses/
│   └── schemas/
├── scripts/                 # Scripts SQL e utilitários
├── tests/                   # Testes automatizados
│   ├── integration/        # Testes de integração
│   └── unit/              # Testes unitários
└── coverage/               # Relatórios de cobertura de testes
```

## ✨ Funcionalidades

- **Autenticação**
  - Registro de usuários
  - Login com JWT
  - Proteção de rotas

- **Gerenciamento de Despesas**
  - CRUD completo de despesas
  - Filtros por período e categoria
  - Paginação de resultados
  - Validação de dados

- **Categorização**
  - Mantimentos
  - Lazer
  - Eletrônica
  - Utilitários
  - Roupas
  - Saúde
  - Outros

## 🛠️ Requisitos

- Go 1.23 ou superior
- PostgreSQL 15+
- Docker & Docker Compose
- Make

## 🚀 Como Executar

### Usando Docker

```bash
# Construir e iniciar os containers
docker-compose up -d

# Verificar logs
docker-compose logs -f
```

### Localmente

1. Clone o repositório
```bash
git clone https://github.com/seu-usuario/ExpenseApi.git
cd ExpenseApi
```

2. Configure as variáveis de ambiente
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. Prepare o banco de dados
```bash
docker-compose up -d postgres
docker-compose exec postgres psql -U expense_user -d postgres -f /docker-entrypoint-initdb.d/init.sql
```

4. Execute os testes
```bash
make test          # Executa todos os testes (unitários e integração)
make test-unit     # Executa apenas testes unitários
make test-integration  # Executa apenas testes de integração
```

5. Inicie o servidor
```bash
make run
```

## 📚 Documentação da API

A documentação completa da API está disponível em:
- Swagger UI: http://localhost:8081/docs
- OpenAPI Spec: `/docs/api.yaml`

### Endpoints Principais

#### Autenticação
- `POST /api/v1/auth/register` - Registro de novo usuário
- `POST /api/v1/auth/login` - Login de usuário

#### Despesas
- `GET /api/v1/expenses` - Lista todas as despesas
- `POST /api/v1/expenses` - Cria uma nova despesa
- `GET /api/v1/expenses/{id}` - Obtém uma despesa específica
- `PUT /api/v1/expenses/{id}` - Atualiza uma despesa
- `DELETE /api/v1/expenses/{id}` - Remove uma despesa

## 🧪 Testes

O projeto inclui testes unitários e de integração:

```bash
# Executar todos os testes (inclui criação do banco de teste)
make test

# Executar apenas testes unitários
make test-unit

# Executar apenas testes de integração
make test-integration

# Executar testes com cobertura
make test-coverage

# Limpar arquivos de cobertura
make clean
```

## 📈 Métricas de Qualidade

- Cobertura de testes: > 80%
- Lint: Golangci-lint
- Documentação: 100% dos endpoints documentados

## 🤝 Contribuindo

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes. 