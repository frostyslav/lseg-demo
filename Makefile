.DEFAULT_GOAL = help
.PHONY: help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-minikube: ## install minikube
ifeq (,$(wildcard /usr/local/bin/minikube))
	sudo curl -Lo /usr/local/bin/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
	sudo chmod +x /usr/local/bin/minikube
endif

start-minikube: ## start minikube
	minikube start

install-registry: ## install registry
	minikube addons enable registry
	ss -tlpn | grep -q 5000 || kubectl port-forward --namespace kube-system $$(kubectl get pods --namespace kube-system --selector actual-registry=true --no-headers -o custom-columns=NAME:.metadata.name) --address 0.0.0.0 5000:5000 &
	docker ps | grep -q socat || docker run -d --name=socat --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$$(minikube ip):5000"

install-knative: ## install knative
	kubectl apply -f https://github.com/knative/serving/releases/download/v0.22.0/serving-crds.yaml
	kubectl apply -f https://github.com/knative/serving/releases/download/v0.22.0/serving-core.yaml
	# kubectl apply -f https://github.com/knative/net-istio/releases/download/v0.22.0/istio.yaml
	# kubectl apply -f https://github.com/knative/net-istio/releases/download/v0.22.0/net-istio.yaml
	kubectl apply -f https://github.com/knative/net-kourier/releases/download/v0.22.0/kourier.yaml
	kubectl patch configmap/config-network --namespace knative-serving --type merge --patch '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
	kubectl patch configmap/config-deployment --namespace knative-serving --type merge --patch '{"data":{"registriesSkippingTagResolving":"localhost:5000"}}'

build-dockerfiles: ## build dockerfiles
	docker build -f Dockerfile.backend -t localhost:5000/knfunc .
	docker push localhost:5000/knfunc
	docker build -f Dockerfile.frontend -t localhost:5000/web-ui .
	docker push localhost:5000/web-ui

install-frontend: ## install frontend
	kubectl apply -f k8s/frontend.yml

install-backend: ## install backend
	kubectl apply -f k8s/backend.yml

destroy: ## destroys minikube
	minikube stop
	minikube delete