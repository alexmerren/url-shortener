GO := go
POSTGRES := psql
DOCKER := docker

EXEC_NAME := url-shortner

INTERNAL_DIR := $(CURDIR)/internal
CMD_DIR := $(CURDIR)/cmd
BIN_DIR := $(CURDIR)/bin
BUILD_DIR := $(CURDIR)/build

DB_USER := postgres
DB_HOST := pgdb
DB_PORT := 5432
DB_INIT := $(BUILD_DIR)/db-init.sql
DB_WIPE := $(BUILD_DIR)/db-wipe.sql

# Application Commands

.PHONY: build
build:
	$(GO) build -o $(BIN_DIR)/... -mod=vendor $(CMD_DIR)/$(EXEC_NAME) 

.PHONY: run
run:
	$(GO) run $(CMD_DIR)/$(EXEC_NAME)

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

# Database Commands

.PHONY: db-init
db-init:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -f $(DB_INIT)

.PHONY: db-wipe
db-wipe:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -f $(DB_WIPE)

# Docker Commands

.PHONY: docker-build
docker-build:
	echo "Not implemented yet..."

.PHONY: docker-run
docker-run:
	echo "Not implemented yet..."
