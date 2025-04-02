build:
    @scripts/build.sh

run: build
    echo "Installing local Helm chart..."
    helm install chakchat k8s --values k8s/values-local.yaml