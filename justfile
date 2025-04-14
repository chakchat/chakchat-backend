build:
    @./build-docker.sh local

run: build
    @echo "Installing local Helm chart..."
    @if [[ $(kubectl config current-context) != 'minikube' ]]; then \
        echo "ERROR: You should set minikube context: 'kubectl config use-context minikube'"; \
        exit 1; \
    fi
    helm install chakchat k8s --values k8s/values-local.yaml
    minikube tunnel

stop:
    helm uninstall chakchat
