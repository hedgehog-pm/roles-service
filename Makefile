.PHONY: fmt lint generate generate-sql migrate docs docs-build

ifeq ($(OS),Windows_NT)
DOCS_PYTHON := python
DOCS_VENV_PYTHON := .venv-docs/Scripts/python.exe
DOCS_CREATE_VENV := if not exist "$(DOCS_VENV_PYTHON)" $(DOCS_PYTHON) -m venv .venv-docs
else
DOCS_PYTHON := python3
DOCS_VENV_PYTHON := .venv-docs/bin/python
DOCS_CREATE_VENV := test -x $(DOCS_VENV_PYTHON) || $(DOCS_PYTHON) -m venv .venv-docs
endif

export PSQL_DSN ?= $(DATABASE_URL)

lint:
	go tool golangci-lint config verify
	go tool golangci-lint run

fmt:
	go tool golangci-lint fmt

.PHONY: fix
fix:
	go tool golangci-lint config verify
	go tool golangci-lint run --fix ./...

generate: generate-sql generate-api

generate-sql:
	@test -n "$$PSQL_DSN" || (echo 'PSQL_DSN must point to a migrated PostgreSQL database'; exit 1)
	@set -eu; \
		trap 'rm -f internal/generated/bob/role/queries.sql internal/generated/bob/roletemplate/queries.sql internal/generated/bob/roleassignment/queries.sql' EXIT; \
		for part in role roletemplate roleassignment; do \
			mkdir -p internal/generated/bob/$$part; \
			cp internal/repository/$$part/queries.sql internal/generated/bob/$$part/queries.sql; \
		done; \
		go tool bobgen-psql -c bobgen.yaml

migrate:
	go run ./cmd/migrate

.PHONY: generate-api
generate-api:
	go tool oapi-codegen -config api/server/roles/oapi-codegen.yaml api/server/roles/openapi.yaml

docs:
	$(DOCS_CREATE_VENV)
	$(DOCS_VENV_PYTHON) -m pip install -r requirements-docs.txt
	$(DOCS_VENV_PYTHON) -m mkdocs serve

docs-build:
	$(DOCS_CREATE_VENV)
	$(DOCS_VENV_PYTHON) -m pip install -r requirements-docs.txt
	$(DOCS_VENV_PYTHON) -m mkdocs build --strict
