IMAGE_NAME := transaction-temporal-workflow/api-protoc-go
VERSION := v0.1.0
DOCKERFILE := protoc.Dockerfile

TMP_DIR := ./tmp

DOCKER_IMAGE_ID := $(shell docker images -q $(IMAGE_NAME):$(VERSION))

define api_protoc_go
	docker run --rm -v ${PWD}:/generate \
		$(IMAGE_NAME):$(VERSION) \
		-c \
		"protoc \
			--go_out=plugins=grpc:$(TMP_DIR) \
			./api/*.proto"
endef

.PHONY: $(DOCKERFILE)
$(DOCKERFILE):
	@if [ "$(DOCKER_IMAGE_ID)" = "" ]; then \
		docker buildx build -f $(DOCKERFILE) -t $(IMAGE_NAME):$(VERSION) --output=type=docker .; \
	fi;

$(TMP_DIR):
	@if [ ! -d $(TMP_DIR) ]; then mkdir tmp; fi

.PHONY: generate
generate: $(TMP_DIR) $(DOCKERFILE)
ifeq ($(package),api)
	@echo "generating api"
	$(api_protoc_go)

	@cp -r $(TMP_DIR)/api/* ./api 2>/dev/null || :
	@sudo rm -rf $(TMP_DIR)
endif
