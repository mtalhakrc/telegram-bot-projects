
MIGRATION_NAME := $(shell date  +%Y%m%d%H%M%S)_$(name)
.PHONY: migrationdir
migrationdir: migrationfile
	mkdir -p common/migrations/$(MIGRATION_NAME)

.PHONY: migrationfile
migrationfile:
	cp common/migrations/template.txt common/migrations/$(MIGRATION_NAME).go

.PHONY: migrate
migrate:
	go run common/cmd/main.go -env=dev init
	go run common/cmd/main.go -env=dev migrate
	go run common/cmd/main.go -env=dev status

.PHONY: rollback
rollback:
	go run common/cmd/main.go -env=dev rollback
