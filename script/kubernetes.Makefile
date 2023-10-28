GO_APP_KUBERNETES_DEPLOYMENT := app-deployment.yml
GO_APP_KUBERNETES_SERVICE := app-service.yml


.PHONY: deploy-argocd
deploy-argocd:
	@if [ "$(shell kubectl get namespaces | grep argocd)" = "" ]; then \
		kubectl create namespace argocd; \
	fi;

	@kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

#	@kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'

#	@if [ "$(shell sudo netstat -tnlp | grep :8080)" = "" ]; then \
#		@kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 & \
#	fi;

	@kubectl port-forward svc/argocd-server -n argocd 8080:443 >/dev/null 2>&1 & \

	@argocd login localhost:8080 --username admin --password admin12345 --insecure

#	@argocd cluster add -y $(kubectl config get-contexts -o name)

.PHONY: deploy-kube
deploy-kube:
	@echo "deploy kubernetes resource $(APP_NAME)"
	@if [ ! -d ./build ]; then mkdir ./build; fi

	@cp -r ../../../$(GO_APP_KUBERNETES_DEPLOYMENT) ./build/deployment.yml
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/deployment.yml
	@sed -i'' -e 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/deployment.yml

	@cp -r ../../../$(GO_APP_KUBERNETES_SERVICE) ./build/service.yml
	@sed -i'' -e 's#appname#$(APP_NAME)#g' ./build/service.yml
	@sed -i'' -e 's#appimagename#$(APP_IMAGE_NAME)#g' ./build/service.yml

	@chmod +x ../../../script/convert-dotenv-to-kube-env.sh && \
		../../../script/convert-dotenv-to-kube-env.sh .env ./build/deployment.yml

	@minikube image load $(APP_IMAGE_NAME)

	@kubectl apply -f ./build/deployment.yml
	@kubectl apply -f ./build/service.yml

.PHONY: deploy-kube-dependency
deploy-kube-dependency: deploy-kube-postgres deploy-kube-redis

.PHONY: deploy-kube-postgres
deploy-kube-postgres:
	@echo "deploy postgres on kubernetes"
	@kubectl apply -f ./postgres-volume.yml

	@helm repo add --force-update bitnami https://charts.bitnami.com/bitnami
	@helm upgrade --install postgresql bitnami/postgresql -f ./postgres-config.yml --output=json | \
		jq -r "[.name, .info.description] | @sh" | xargs printf "%s %s\n"

.PHONY: deploy-kube-redis
deploy-kube-redis:
	@echo "deploy redis on kubernetes"
	@kubectl apply -f ./redis-volume.yml

	@helm repo add --force-update bitnami https://charts.bitnami.com/bitnami
	@helm upgrade --install redis bitnami/redis -f ./redis-config.yml --output=json | \
		jq -r "[.name, .info.description] | @sh" | xargs printf "%s %s\n"

.PHONY: convert-dotenv-to-kube-env
convert-dotenv-to-kube-env:
	@chmod +x ./convert-dotenv-to-kube-env.sh && ./convert-dotenv-to-kube-env.sh
