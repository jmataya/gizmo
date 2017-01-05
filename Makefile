FLYWAY=flyway -configFile=sql/flyway.conf -locations=filesystem:sql/

migrate:
	$(FLYWAY) migrate

reset:
	dropdb --if-exists gizmo
	dropuser --if-exists gizmo
	createuser -s gizmo
	createdb gizmo

.PHONY: migrate reset