TMP_DIR := ./tmp
MIGRATION_DIR := $(PWD)/migration
TRANSACTION_SERVER_DIR := ./cmd/server/transaction
TRANSACTION_PUBSUB_DIR := ./cmd/pubsub/transaction
TRANSACTION_CRON_DIR := ./cmd/cron/transaction
TRANSACTION_WORKER_DIR := ./cmd/worker/transaction

PROTOC_DOCKERFILE := protoc.Dockerfile
PROTOC_IMAGE_NAME := transaction-temporal-workflow/api-protoc-go:latest
PROTOC_IMAGE_ID := $(shell docker images -q $(PROTOC_IMAGE_NAME))

define api_protoc_go
	docker run --rm -v ${PWD}:/generate \
		$(PROTOC_IMAGE_NAME) \
		-c \
		"protoc \
			--go_out=plugins=grpc:$(TMP_DIR) \
			./api/*.proto"
endef

define api_protoc_nodejs
	docker run --rm -v ${PWD}:/generate \
		$(PROTOC_IMAGE_NAME) \
		-c \
			'grpc_tools_node_protoc \
				--plugin="$$(which protoc-gen-es)" \
				--es_out $(TMP_DIR) \
				--es_opt target=ts \
				--plugin="$$(which protoc-gen-connect-es)" \
				--connect-es_out $(TMP_DIR) \
				--connect-es_opt target=ts \
				--plugin="$$(which protoc-gen-ts_proto)" \
				--ts_proto_out=$(TMP_DIR) \
				--ts_proto_opt=outputServices=grpc-js \
				--ts_proto_opt=env=node \
				--ts_proto_opt=esModuleInterop=true \
				./api/*.proto'
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
	$(api_protoc_nodejs)

	@cp -r $(TMP_DIR)/api/*.go ./api 2>/dev/null || :
	@cp -r $(TMP_DIR)/api/*.ts ./gateway-nodejs/api 2>/dev/null || :
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

.PHONY: deploy-argocd
deploy-argocd:
	@if [ "$(shell kubectl get namespaces | grep argocd)" = "" ]; then \
		@kubectl create namespace argocd; \
	fi;

	@kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

#	@kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'

#	@if [ "$(shell sudo netstat -tnlp | grep :8080)" = "" ]; then \
#		@kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 & \
#	fi;

	@kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 & \

	@argocd login localhost:8080 --username admin --password admin12345 --insecure

#	@argocd cluster add -y $(kubectl config get-contexts -o name)
