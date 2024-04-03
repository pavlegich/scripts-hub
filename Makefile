DATABASE_DSN = postgresql://localhost:5432/postgres

DOC_PORT = 6060

SERVER_BINARY_NAME = server
SERVER_PACKAGE_PATH = ./cmd/server
SERVER_ADDR = localhost:8080

# ====================
# HELPERS
# ====================

## help: show this help message
help:
	@echo
	@echo 'usage: make <target>'
	@echo
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
	@echo

# ====================
# QUALITY
# ====================

## tidy: format code and tidy mod file
tidy:
	go fmt ./...
	go mod tidy -v

# ====================
# DEVELOPMENT
# ====================

## build: build the server
build:
	go build -o /tmp/bin/$(SERVER_BINARY_NAME) $(SERVER_PACKAGE_PATH)

## run: run the server
run: build
	/tmp/bin/$(SERVER_BINARY_NAME) -a=$(SERVER_ADDR) -d=$(DATABASE_DSN)

# ====================
# DOCUMENTATION
# ====================

## doc: generate documentation on http port
doc:
	@echo 'open http://localhost:$(DOC_PORT)/pkg/github.com/pavlegich/scripts-hub/?m=all'
	godoc -http=:$(DOC_PORT)

.PHONY: help tidy build run doc
