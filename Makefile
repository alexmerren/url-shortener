GO = go

NAME = url-shortner

INTERNAL_DIR = $(CURDIR)/internal
CMD_DIR = $(CURDIR)/cmd
BIN_DIR = $(CURDIR)/bin

.PHONY: build
build:
	@echo "Building executable into $(CMD_DIR)/$(NAME)"
	@go build -o $(BIN_DIR)/$(NAME) -mod=vendor $(CMD_DIR)/$(NAME) 

.PHONY: run
run:
	@$(BIN_DIR)/$(NAME)

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: db-init
db-init:
	@sqlite3 dev/database.db < dev/schema.sql
