GO := go
POSTGRES := psql
DOCKERCOMPOSE := docker-compose

INTERNAL_DIR := $(CURDIR)/internal
CMD_DIR := $(CURDIR)/cmd
DIST_DIR := $(CURDIR)/dist
BUILD_DIR := $(CURDIR)/build

DB_USER := postgres
DB_HOST := localhost
DB_PORT := 5432
DB_NAME := urlshortener

DB_INIT := $(BUILD_DIR)/db-init.sql
DB_CREATE := $(BUILD_DIR)/db-create.sql
DB_WIPE := $(BUILD_DIR)/db-wipe.sql

# Application Commands

.PHONY: build
build:
	$(GO) build -o $(DIST_DIR)/ -mod=vendor $(CMD_DIR)/...

.PHONY: run
run:
	$(GO) run $(CMD_DIR)/url-shortener

.PHONY: vendor
vendor:
	$(GO) mod tidy
	$(GO) mod vendor

# Database Commands

.PHONY: db-create
db-create:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -f $(DB_CREATE)

.PHONY: db-init
db-init:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -d $(DB_NAME) -f $(DB_INIT)

.PHONY: db-wipe
db-wipe:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -d $(DB_NAME) -f $(DB_WIPE)

.PHONY: db-exec
db-exec:
	$(POSTGRES) -h $(DB_HOST) -p $(DB_PORT)  -U $(DB_USER) -d $(DB_NAME)

# Docker Commands

.PHONY: docker-build
docker-build:
	$(DOCKERCOMPOSE) build --no-cache

.PHONY: docker-run
docker-run:
	$(DOCKERCOMPOSE) up -d
