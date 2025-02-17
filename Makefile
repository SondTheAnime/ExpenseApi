.PHONY: test test-unit test-integration test-coverage create-test-db

# Variáveis
COVERAGE_DIR=coverage
COVERAGE_FILE=$(COVERAGE_DIR)/coverage.out
COVERAGE_HTML=$(COVERAGE_DIR)/coverage.html

# Comandos principais
test: create-test-db test-unit test-integration

test-unit:
	@echo "Executando testes unitários..."
	@go test -v ./tests/unit/...

test-integration:
	@echo "Executando testes de integração..."
	@go test -v ./tests/integration/...

test-coverage:
	@echo "Gerando relatório de cobertura..."
	@mkdir -p $(COVERAGE_DIR)
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Relatório de cobertura gerado em $(COVERAGE_HTML)"

create-test-db:
	@echo "Criando banco de dados de teste..."
	@docker compose exec -T postgres psql -U expense_user -d postgres -c "DROP DATABASE IF EXISTS expense_test_db;"
	@docker compose exec -T postgres psql -U expense_user -d postgres -f /docker-entrypoint-initdb.d/init.sql
	@docker compose exec -T postgres psql -U expense_user -d postgres -c "CREATE DATABASE expense_test_db WITH TEMPLATE expense_db;"

# Comandos de limpeza
clean:
	@rm -rf $(COVERAGE_DIR)
	@echo "Diretório de cobertura removido" 