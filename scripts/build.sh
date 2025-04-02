#!/bin/bash


debug=false
for arg in "$@"; do
    if [ "$arg" == "--debug" ]; then
            debug=true
    fi
done

images=("identity" "file-storage" "messaging" "user")
paths=("identity-service/" "file-storage-service/" "messaging-service/" "user-service/")
tag=local

build-parallel() {
    for ((i in "${#images[@]}"; do
        minikube image build "${images[$image]}" -t "$image:$tag" &
    done

    if ! wait; then
        echo "At least one build failed! Run it in debug mode"
        exit 1
    fi
}

build-debug() {
    echo "Bulding docker images..."
    for image in "${!images[@]}"; do
        if ! minikube image build "${images[$image]}" -t "$image:$tag"; then
            echo "Building $image failed"
        fi
    done
}

if $debug; then 
    build-debug
else
    build-parallel
fi

echo "Successfully built docker images!"