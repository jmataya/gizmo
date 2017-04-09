FLYWAY=flyway -configFile=sql/flyway.conf -locations=filesystem:sql/
FLYWAY_TEST=flyway -configFile=sql/flyway.test.conf -locations=filesystem:sql/

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

migrate-test:
	$(FLYWAY_TEST) migrate

reset:
	dropdb --if-exists gizmo
	dropuser --if-exists gizmo
	createuser -s gizmo
	createdb gizmo

reset-test:
	dropdb --if-exists gizmo_test
	createdb gizmo_test
	@make migrate-test

test: info-test
	go test -p 1 . ./models

.PHONY: glide info-test migrate migrate-test reset reset-test test
