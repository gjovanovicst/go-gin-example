.PHONY: build run clean migrate-up migrate-down migrate-version migrate-to help seed-dev seed-prod seed-rollback seed-status seed-list new-seed

# Build the application
build:
	go build -o bin/go-gin-example main.go
	go build -o bin/migrate cmd/migrate/main.go
	go build -o bin/seed cmd/seed/main.go

# Run the application
run: build
	./bin/go-gin-example

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f go-gin-example.exe
	rm -f tmp/main.exe

# Run all pending migrations
migrate-up:
	go run cmd/migrate/main.go -action=up

# Rollback the last migration
migrate-down:
	go run cmd/migrate/main.go -action=down

# Show current migration version
migrate-version:
	go run cmd/migrate/main.go -action=version

# Migrate to a specific version (use: make migrate-to VERSION=1)
migrate-to:
	@if [ -z "$(VERSION)" ]; then \
		echo "Usage: make migrate-to VERSION=<version_number>"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -action=migrate -version=$(VERSION)

# Create a new migration file (use: make new-migration NAME=create_users_table)
new-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "Usage: make new-migration NAME=<migration_name>"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${TIMESTAMP}_$(NAME).up.sql; \
	touch migrations/$${TIMESTAMP}_$(NAME).down.sql; \
	echo "Created migration files:"; \
	echo "  migrations/$${TIMESTAMP}_$(NAME).up.sql"; \
	echo "  migrations/$${TIMESTAMP}_$(NAME).down.sql"

# Create a new seed file (use: make new-seed ENV=development NAME=create_sample_data)
new-seed:
	@if [ -z "$(ENV)" ] || [ -z "$(NAME)" ]; then \
		echo "Usage: make new-seed ENV=<environment> NAME=<seed_name>"; \
		echo "Example: make new-seed ENV=development NAME=create_sample_articles"; \
		echo "Example: make new-seed ENV=production NAME=create_admin_user"; \
		exit 1; \
	fi
	./scripts/new-seed.sh $(ENV) $(NAME)

# Install golang-migrate CLI tool
install-migrate-cli:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Reset database (WARNING: This will drop all data!)
migrate-reset:
	@echo "WARNING: This will drop all data from the database!"
	@echo "Press Ctrl+C to cancel, or Enter to continue..."
	@read
	go run cmd/migrate/main.go -action=migrate -version=0
	go run cmd/migrate/main.go -action=up
	@echo "Database reset complete. You can now run seeds:"
	@echo "  make seed-dev    - for development environment"
	@echo "  make seed-prod   - for production environment"

# Seeding commands (tracked - recommended)
seed-dev:
	go run cmd/seed/main.go -action=run -env=development

seed-prod:
	go run cmd/seed/main.go -action=run -env=production

seed-rollback:
	@if [ -z "$(ENV)" ]; then \
		echo "Usage: make seed-rollback ENV=<environment>"; \
		echo "Example: make seed-rollback ENV=development"; \
		exit 1; \
	fi
	go run cmd/seed/main.go -action=rollback -env=$(ENV)

seed-status:
	@if [ -z "$(ENV)" ]; then \
		echo "Checking status for all environments..."; \
		go run cmd/seed/main.go -action=status -env=development; \
		go run cmd/seed/main.go -action=status -env=production; \
	else \
		go run cmd/seed/main.go -action=status -env=$(ENV); \
	fi

seed-list:
	go run cmd/seed/main.go -action=list

# Detailed seed information
seed-info:
	@if [ -z "$(ENV)" ]; then \
		echo "Usage: make seed-info ENV=<environment>"; \
		echo "Example: make seed-info ENV=development"; \
		exit 1; \
	fi
	go run cmd/seed/main.go -action=detailed-status -env=$(ENV)

# Manual seeding (no tracking)
seed-manual:
	@if [ -z "$(ENV)" ]; then \
		echo "Usage: make seed-manual ENV=<environment>"; \
		echo "Example: make seed-manual ENV=development"; \
		exit 1; \
	fi
	go run cmd/seed/main.go -action=run -env=$(ENV) -manual

# Setup full environment (migrations + seeds)
setup-dev: migrate-up seed-dev
	@echo "Development environment setup complete!"

setup-prod: migrate-up seed-prod
	@echo "Production environment setup complete!"

# Show available commands
help:
	@echo "Available commands:"
	@echo ""
	@echo "Application:"
	@echo "  build             - Build the application and migration tools"
	@echo "  run               - Build and run the application"
	@echo "  clean             - Clean build artifacts"
	@echo ""
	@echo "Migrations:"
	@echo "  migrate-up        - Run all pending migrations"
	@echo "  migrate-down      - Rollback the last migration"
	@echo "  migrate-version   - Show current migration version"
	@echo "  migrate-to        - Migrate to a specific version (use VERSION=N)"
	@echo "  new-migration     - Create new migration files (use NAME=migration_name)"
	@echo "  new-seed          - Create new seed files (use ENV=environment NAME=seed_name)"
	@echo "  migrate-reset     - Reset database (WARNING: drops all data!)"
	@echo ""
	@echo "Seeding:"
	@echo "  seed-dev          - Run development seeds"
	@echo "  seed-prod         - Run production seeds"
	@echo "  seed-rollback     - Rollback seeds (use ENV=environment)"
	@echo "  seed-status       - Show seed status (optional ENV=environment)"
	@echo "  seed-info         - Show detailed seed info (use ENV=environment)"
	@echo "  seed-list         - List available seed environments"
	@echo "  seed-manual       - Manual seeding without tracking (use ENV=environment)"
	@echo ""
	@echo "Environment Setup:"
	@echo "  setup-dev         - Full setup for development (migrations + dev seeds)"
	@echo "  setup-prod        - Full setup for production (migrations + prod seeds)"
	@echo ""
	@echo "Other:"
	@echo "  install-migrate-cli - Install golang-migrate CLI tool"
	@echo "  help              - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make migrate-up"
	@echo "  make migrate-to VERSION=1"
	@echo "  make new-migration NAME=add_user_email_column"
	@echo "  make new-seed ENV=development NAME=create_sample_articles"
	@echo "  make seed-dev"
	@echo "  make seed-info ENV=development"
	@echo "  make seed-rollback ENV=development"
	@echo "  make setup-dev"
