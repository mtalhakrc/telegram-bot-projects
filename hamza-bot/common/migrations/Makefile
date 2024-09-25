.PHONY: migrationdir
MIGRATION_NAME := $(shell date  +%Y%m%d%H%M%S)_$(name)
migrationdir: migrationfile
	mkdir -p common/migrations/$(MIGRATION_NAME)

.PHONY: migrationfile
migrationfile:
	cp common/migrations/template.txt common/migrations/$(MIGRATION_NAME).go

.PHONY: migrate
migrate:
	go run ./cmd/main.go init
	go run ./cmd/main.go migrate
	go run ./cmd/main.go status

.PHONY: rollback
rollback:
	go run pkg/cmd/main.go rollback

.PHONY: seed
seed:
	go run pkg/cmd/main.go seed

.PHONY: resetdb
resetdb: rollback migrate seed