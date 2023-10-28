include ./script/database.Makefile
include ./script/proto.Makefile
include ./script/docker.Makefile
include ./script/kubernetes.Makefile

.PHONY: build-image-all
build-image-all:
	@echo "build all app images"
	@make -C cmd/server/transaction build-image
#	@make -C cmd/worker/transaction build-image
#	@make -C cmd/worker/user build-image
#	@make -C cmd/pubsub/transaction build-image
#	@make -C cmd/pubsub/user build-image
#	@make -C cmd/cron/transaction build-image
#	@make -C cmd/cron/user build-image

.PHONY: remove-image-all
remove-image-all:
	@echo "remove all app images"
	@make -C cmd/server/transaction remove-image
	@make -C cmd/worker/transaction remove-image
	@make -C cmd/worker/user remove-image
	@make -C cmd/pubsub/transaction remove-image
	@make -C cmd/pubsub/user remove-image
	@make -C cmd/cron/transaction remove-image
	@make -C cmd/cron/user remove-image

.PHONY: deploy-kube-all
deploy-kube-all: deploy-kube-dependency
	@echo "deploy all app to kubernetes"
	@make -C cmd/server/transaction deploy-kube
#	@make -C cmd/worker/transaction build-image
#	@make -C cmd/worker/user build-image
#	@make -C cmd/pubsub/transaction build-image
#	@make -C cmd/pubsub/user build-image
#	@make -C cmd/cron/transaction build-image
#	@make -C cmd/cron/user build-image
