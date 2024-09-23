GO_APP_KUBERNETES_DEPLOYMENT := app-deployment.yml
GO_APP_KUBERNETES_JOB := app-job.yml
GO_APP_KUBERNETES_SERVICE := app-service.yml

.PHONY: prepare-kube
prepare-kube:
	@minikube start --driver=docker

.PHONY: deploy-kube
deploy-kube:
	@echo "deploy kubernetes resource $(APP_NAME)"
	@if [ ! -d ./build ]; then mkdir ./build; fi

ifneq ($(APP_TYPE),cron)
	@cp -r ../../../script/$(GO_APP_KUBERNETES_DEPLOYMENT) ./build/deployment.yml
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/deployment.yml
	@sed -i'' -e 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/deployment.yml
	@chmod +x ../../../script/convert-dotenv-to-kube-env.sh && \
			../../../script/convert-dotenv-to-kube-env.sh .kubernetes.env ./build/deployment.yml
	@kubectl delete --ignore-not-found=true -f ./build/deployment.yml
else
	@cp -r ../../../script/$(GO_APP_KUBERNETES_JOB) ./build/job.yml
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/job.yml
	@sed -i'' -e 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/job.yml
	@chmod +x ../../../script/convert-dotenv-to-kube-env.sh && \
			../../../script/convert-dotenv-to-kube-env.sh .kubernetes.env ./build/job.yml
	@kubectl delete --ignore-not-found=true -f ./build/job.yml
endif

ifeq ($(APP_TYPE),server)
	@cp -r ../../../script/$(GO_APP_KUBERNETES_SERVICE) ./build/service.yml
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/service.yml
	@sed -i'' -e 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/service.yml

	@kubectl delete --ignore-not-found=true -f ./build/service.yml
endif

	@sleep 5
	@if [ -n "$(shell minikube image ls | grep docker.io/${APP_IMAGE_NAME})" ]; then \
  		minikube image rm docker.io/$(APP_IMAGE_NAME); \
  	fi

	@minikube image load $(APP_IMAGE_NAME)

ifneq ($(APP_TYPE),cron)
	@kubectl apply -f ./build/deployment.yml
else
	@kubectl apply -f ./build/job.yml
endif

ifeq ($(APP_TYPE),server)
	@kubectl apply -f ./build/service.yml
endif

.PHONY: deploy-kube-dependency
deploy-kube-dependency: deploy-kube-postgres deploy-kube-redis deploy-kube-rabbitmq deploy-kube-temporal

.PHONY: deploy-kube-postgres
deploy-kube-postgres:
	@echo "deploy postgres on kubernetes"
	@kubectl apply -f ./script/postgres-volume.yml

	@helm repo add --force-update bitnami https://charts.bitnami.com/bitnami
	@helm upgrade --install postgresql bitnami/postgresql -f ./script/postgres-config.yml --output=json | \
		jq -r "[.name, .info.description] | @sh" | xargs printf "%s %s\n"

.PHONY: deploy-kube-redis
deploy-kube-redis:
	@echo "deploy redis on kubernetes"
	@kubectl apply -f ./script/redis-volume.yml

	@helm repo add --force-update bitnami https://charts.bitnami.com/bitnami
	@helm upgrade --install redis bitnami/redis -f ./script/redis-config.yml --output=json | \
		jq -r "[.name, .info.description] | @sh" | xargs printf "%s %s\n"

.PHONY: deploy-kube-rabbitmq
deploy-kube-rabbitmq:
	@echo "deploy rabbitmq on kubernetes"

	@helm repo add --force-update bitnami https://charts.bitnami.com/bitnami
	@helm upgrade --install rabbitmq bitnami/rabbitmq -f ./script/rabbitmq-config.yml --output=json | \
		jq -r "[.name, .info.description] | @sh" | xargs printf "%s %s\n"

.PHONY: deploy-kube-temporal
deploy-kube-temporal:
	@echo "deploy temporal on kubernetes"

	@cd ./thirdparty/temporal/helm-charts && \
		helm dependencies update && \
		helm upgrade temporal . --install -f ../../../script/temporal-config.yml

	@kubectl exec $(shell kubectl get pods --selector=app.kubernetes.io/component=admintools -o jsonpath='{.items[*].metadata.name}') -- tctl --namespace default namespace register
