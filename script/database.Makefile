MIGRATION_DIR := $(PWD)/migration


.PHONY: migrate-up
migrate-up:
	@echo "migrate up database"
	# TODO: move credential to .env
	docker run -v $(MIGRATION_DIR):/migrations \
		--network host \
		migrate/migrate \
        	-path=/migrations/ \
        	-database postgres://app:app@localhost:5433/app?sslmode=disable \
        	up

.PHONY: migrate-down
migrate-down:
	@echo "migrate down database"
	# TODO: move credential to .env
	docker run -v $(MIGRATION_DIR):/migrations \
		--network host \
		migrate/migrate \
        	-path=/migrations/ \
        	-database postgres://app:app@localhost:5433/app?sslmode=disable \
        	down -all
