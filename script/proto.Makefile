TMP_DIR := ./tmp
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
