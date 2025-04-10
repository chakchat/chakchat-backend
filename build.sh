#!/bin/bash

debug=false
for arg in "$@"; do
    if [ "$arg" == "--debug" ]; then
            debug=true
    fi
done

images=(
    "identity" 
    "file-storage" 
    "messaging" 
    "messaging-pg-migrate" 
    "user" 
    "sms-service-stub"
)
paths=(
    "identity-service/" 
    "file-storage-service/" 
    "messaging-service/" 
    "messaging-service/" 
    "user-service/" 
    "stubs/sms-service-stub/"
)
dockerfiles=(
    "Dockerfile"
    "Dockerfile"
    "Dockerfile"
    "migrate.Dockerfile"
    "Dockerfile"
    "Dockerfile"
)
tag=local

build-parallel() {
    for i in "${!images[@]}"; do
        minikube image build "${paths[$i]}" -f "${dockerfiles[$i]}" -t "${images[$i]}:$tag" &> /dev/null &
    done

    if ! wait; then
        echo "At least one build failed! Run it in debug mode"
        exit 1
    fi
}

build-debug() {
    echo "Bulding docker images..."
    for i in "${!images[@]}"; do
        if ! minikube image build "${paths[$i]}" -f "${dockerfiles[$i]}" -t "${images[$i]}:$tag"; then
            echo "Building ${images[$i]} failed"
            exit 1
        fi
    done
}

if $debug; then 
    build-debug
else
    build-parallel
fi

echo "Successfully built docker images!"