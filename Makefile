TMP_DIR := ./tmp
MIGRATION_DIR := $(PWD)/migration
TRANSACTION_SERVER_DIR := ./cmd/server/transaction
TRANSACTION_PUBSUB_DIR := ./cmd/pubsub/transaction
TRANSACTION_CRON_DIR := ./cmd/cron/transaction
TRANSACTION_WORKER_DIR := ./cmd/worker/transaction

PROTOC_DOCKERFILE := protoc.Dockerfile
PROTOC_IMAGE_NAME := transaction-temporal-workflow/api-protoc-go:v0.1.0
PROTOC_IMAGE_ID := $(shell docker images -q $(PROTOC_IMAGE_NAME))

GO_APP_DOCKERFILE := go-app.Dockerfile
TRANSACTION_SERVER_IMAGE_NAME := transaction-temporal-workflow/transaction-server:v0.1.0
TRANSACTION_SERVER_IMAGE_ID := $(shell docker images -q $(TRANSACTION_SERVER_IMAGE_NAME))
TRANSACTION_SERVER_NAME := transaction-server
TRANSACTION_SERVER_DIR := cmd/server/transaction

TRANSACTION_WORKER_IMAGE_NAME := transaction-temporal-workflow/transaction-worker:v0.1.0
TRANSACTION_WORKER_IMAGE_ID := $(shell docker images -q $(TRANSACTION_WORKER_IMAGE_NAME))
TRANSACTION_WORKER_NAME := transaction-worker
TRANSACTION_WORKER_DIR := cmd/worker/transaction

TRANSACTION_PUBSUB_IMAGE_NAME := transaction-temporal-workflow/transaction-pubsub:v0.1.0
TRANSACTION_PUBSUB_IMAGE_ID := $(shell docker images -q $(TRANSACTION_PUBSUB_IMAGE_NAME))
TRANSACTION_PUBSUB_NAME := transaction-pubsub
TRANSACTION_PUBSUB_DIR := cmd/pubsub/transaction

TRANSACTION_CRON_IMAGE_NAME := transaction-temporal-workflow/transaction-cron:v0.1.0
TRANSACTION_CRON_IMAGE_ID := $(shell docker images -q $(TRANSACTION_CRON_IMAGE_NAME))
TRANSACTION_CRON_NAME := transaction-cron
TRANSACTION_CRON_DIR := cmd/cron/transaction

USER_WORKER_IMAGE_NAME := transaction-temporal-workflow/user-worker:v0.1.0
USER_WORKER_IMAGE_ID := $(shell docker images -q $(USER_WORKER_IMAGE_NAME))
USER_WORKER_NAME := user-worker
USER_WORKER_DIR := cmd/worker/user

USER_PUBSUB_IMAGE_NAME := transaction-temporal-workflow/user-pubsub:v0.1.0
USER_PUBSUB_IMAGE_ID := $(shell docker images -q $(USER_PUBSUB_IMAGE_NAME))
USER_PUBSUB_NAME := user-pubsub
USER_PUBSUB_DIR := cmd/pubsub/user

USER_CRON_IMAGE_NAME := transaction-temporal-workflow/user-cron:v0.1.0
USER_CRON_IMAGE_ID := $(shell docker images -q $(USER_CRON_IMAGE_NAME))
USER_CRON_NAME := user-cron
USER_CRON_DIR := cmd/cron/user

define api_protoc_go
	docker run --rm -v ${PWD}:/generate \
		$(PROTOC_IMAGE_NAME) \
		-c \
		"protoc \
			--go_out=plugins=grpc:$(TMP_DIR) \
			./api/*.proto"
endef

.PHONY: $(PROTOC_DOCKERFILE)
$(PROTOC_DOCKERFILE):
	@if [ "$(PROTOC_IMAGE_ID)" = "" ]; then \
		docker buildx build -f $(PROTOC_DOCKERFILE) -t $(PROTOC_IMAGE_NAME) --output=type=docker .; \
	fi;

$(TMP_DIR):
	@if [ ! -d $(TMP_DIR) ]; then mkdir $(TMP_DIR); fi

.PHONY: generate
generate: $(TMP_DIR) $(PROTOC_DOCKERFILE)
ifeq ($(package),api)
	@echo "generating api"
	$(api_protoc_go)

	@cp -r $(TMP_DIR)/api/* ./api 2>/dev/null || :
	@sudo rm -rf $(TMP_DIR)
endif

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

.PHONY: build-all-images
build-all-images:
	@echo "build all app images"
	@make -C cmd/server/transaction build-image
	@make -C cmd/worker/transaction build-image
	@make -C cmd/worker/user build-image
	@make -C cmd/pubsub/transaction build-image
	@make -C cmd/pubsub/user build-image
	@make -C cmd/cron/transaction build-image
	@make -C cmd/cron/user build-image

.PHONY: remove-all-images
remove-all-images:
	@echo "remove all app images"
	@make -C cmd/server/transaction remove-image
	@make -C cmd/worker/transaction remove-image
	@make -C cmd/worker/user remove-image
	@make -C cmd/pubsub/transaction remove-image
	@make -C cmd/pubsub/user remove-image
	@make -C cmd/cron/transaction remove-image
	@make -C cmd/cron/user remove-image
