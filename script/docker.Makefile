GO_APP_DOCKERFILE := go-app.Dockerfile
APP_IMAGE_ID := $(shell docker images -q $(APP_IMAGE_NAME))


.PHONY: build-image
build-image:
	@echo "build $(APP_NAME)"
	@if [ ! -d ./build ]; then mkdir ./build; fi

	@cp -r ../../../script/$(GO_APP_DOCKERFILE) ./build/Dockerfile
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/Dockerfile
	@sed -i'' -e 's#appdir#$(APP_DIR)#g' ./build/Dockerfile

	@docker buildx build -f ./build/Dockerfile -t $(APP_IMAGE_NAME) --output=type=docker ../../..

.PHONY: remove-image
remove-image:
	@docker image rm $(APP_IMAGE_ID) --force
