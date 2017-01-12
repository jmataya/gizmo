FLYWAY=flyway -configFile=sql/flyway.conf -locations=filesystem:sql/

# Buildkite highlighting
RED = \033[33m
NO_COLOR = \033[0m
baseheader = @echo "---$(1)$(RED)$(2)$(NO_COLOR)"
header = $(call baseheader, $(1), gizmo)

info-test:
	$(call header, Testing)

glide:
	glide install

migrate:
	$(FLYWAY) migrate

reset:
	dropdb --if-exists gizmo
	dropuser --if-exists gizmo
	createuser -s gizmo
	createdb gizmo

test: info-test glide reset migrate
	go run examples/simple.go
	go run examples/sku.go

.PHONY: glide info-test migrate reset test