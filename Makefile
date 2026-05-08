.PHONY: setup up down dev generate test test-watch

# Instala as ferramentas globais na máquina (Air e SQLC)
setup:
	@echo "🛠️ Instalando ferramentas de desenvolvimento..."
	go install github.com/air-verse/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "✅ Ferramentas instaladas com sucesso!"

# Sobe o banco de dados via Docker
up:
	@echo "🐳 Subindo o banco de dados..."
	docker-compose up -d

# Derruba o banco de dados
down:
	@echo "🛑 Parando o banco de dados..."
	docker-compose down

# Roda o gerador de código do SQLC
generate:
	@echo "⚙️ Gerando código Go a partir do SQL..."
	sqlc generate
	@echo "✅ Código gerado em internal/db!"

# O nosso comando principal para trabalhar!
# Ele garante que o banco suba antes de rodar o Air.
dev: up
	@echo "🚀 Iniciando o servidor em modo de desenvolvimento..."
	air

# Roda todos os testes do projeto uma vez
pkg ?= ...

test:
	go test -v ./$(pkg)

# Modo TDD: Observa mudanças e roda os testes automaticamente via Air

# Modo TDD: Agora aceita o argumento 'pkg'
# Exemplo: make test-watch pkg=faq/

test-watch:
	@# Se pkg for "...", o Air usa a config padrão.
	@# Se pkg for algo específico, sobrescrevemos o comando do Air via variável de ambiente.
	TEST_PKG=./$(pkg) air -c .air.test.conf