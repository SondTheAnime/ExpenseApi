#!/bin/bash

# Configura o ambiente de teste
export DB_NAME="expense_test_db"
export DB_USER="postgres"
export DB_PASSWORD="postgres"
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_SSL_MODE="disable"
export JWT_SECRET="test_secret"
export JWT_EXPIRES_IN="3600"
export JWT_REFRESH_TOKEN="604800"

# Cria o banco de dados de teste
psql -U postgres -h localhost -f scripts/create_test_db.sql

# Executa os testes
go test -v ./tests/unit/... ./tests/integration/... 