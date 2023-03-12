APP_IMAGE_NAME := transaction-temporal-workflow/transaction-server:latest
APP_NAME := transaction-server
APP_DIR := cmd/server/transaction

GO_APP_DOCKERFILE := go-app.Dockerfile
GO_APP_KUBERNETES_DEPLOYMENT := app-deployment.yml
GO_APP_KUBERNETES_SERVICE := app-service.yml
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

.PHONY: deploy-kube
deploy-kube:
	@echo "deploy kubernetes resource $(APP_NAME)"
	@if [ ! -d ./build ]; then mkdir ./build; fi

	@cp -r ../../../$(GO_APP_KUBERNETES_DEPLOYMENT) ./build/deployment.yml
	@sed -i 's#appname#$(APP_NAME)#g' ./build/deployment.yml
	@sed -i 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/deployment.yml

	@cp -r ../../../$(GO_APP_KUBERNETES_SERVICE) ./build/service.yml
	@sed -i 's#appname#$(APP_NAME)#g' ./build/service.yml
	@sed -i 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/service.yml
