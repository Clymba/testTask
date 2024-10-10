.PHONY: up
up:
	docker compose up -d

.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: migrate
migrate:
	sql-migrate up

.PHONY: install-deps
install-deps:
	go mod download

.PHONY: install-tools
install-tools:
	go get github.com/rubenv/sql-migrate/sql-migrate@latest
	go install github.com/rubenv/sql-migrate/sql-migrate