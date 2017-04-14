FLYWAY=flyway -configFile=sql/flyway.conf -locations=filesystem:sql/
FLYWAY_TEST=flyway -configFile=sql/flyway.test.conf -locations=filesystem:sql/

# Buildkite highlighting
RED = \033[33m
NO_COLOR = \033[0m
baseheader = @echo "---$(1)$(RED)$(2)$(NO_COLOR)"
header = $(call baseheader, $(1), gizmo)

docs:
	$(call header, Starting Docs Server)
	godoc -http=:6060 &
	echo 'View documentation: http://localhost:6060/pkg/github.com/FoxComm/gizmo'

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
	dropdb --if-exists gizmo_test
	dropuser --if-exists gizmo
	createuser -s gizmo
	createdb gizmo
	createdb gizmo_test
	@make migrate

reset-test:
	dropdb --if-exists gizmo_test
	createdb gizmo_test
	@make migrate-test

test: info-test
	go test -p 1 . ./models

.PHONY: glide docs info-test migrate migrate-test reset reset-test test
