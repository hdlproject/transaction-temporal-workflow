APP_IMAGE_NAME := transaction-temporal-workflow/user-cron:latest
APP_NAME := user-cron
APP_DIR := cmd/cron/user

GO_APP_DOCKERFILE := go-app.Dockerfile
APP_IMAGE_ID := $(shell docker images -q $(APP_IMAGE_NAME))

.PHONY: build-image
build-image:
	@echo "build $(APP_NAME)"
	@if [ ! -d ./build ]; then mkdir ./build; fi

	@cp -r ../../../$(GO_APP_DOCKERFILE) ./build/Dockerfile
	@sed -i 's#appname#$(APP_NAME)#g' ./build/Dockerfile
	@sed -i 's#appdir#$(APP_DIR)#g' ./build/Dockerfile

	@if [ "$(APP_IMAGE_ID)" = "" ]; then \
		docker buildx build -f ./build/Dockerfile -t $(APP_IMAGE_NAME) --output=type=docker ../../..; \
	fi;

.PHONY: remove-image
remove-image:
	@if [ "$(APP_IMAGE_ID)" != "" ]; then \
		docker image rm $(APP_IMAGE_ID) --force; \
	fi;
