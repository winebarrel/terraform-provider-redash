.PHONY: redash-setup
redash-setup:
	psql -U postgres -h localhost -p 15432 -f _etc/redash.sql

.PHONY: redash-create-db
redash-create-db:
	docker compose run --rm server create_db
