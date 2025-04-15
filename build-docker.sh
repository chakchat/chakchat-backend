#!/bin/bash

environment=$1
if [ "$environment" != 'test' ] && [ "$environment" != 'prod' ] && [ "$environment" != 'local' ]; then
    echo "You should pass environment argument. Possible values: prod, test, local"
fi

debug=false
nopush=false
for arg in "$@"; do
    if [ "$arg" == "--debug" ]; then
        debug=true
    elif [ "$arg" == "--no-push" ]; then
        nopush=true
    fi
done

if [[ $environment = 'local' ]]; then
    tag=local
else 
    tag=$environment-$(git rev-parse --short HEAD  )
fi

repository=""
if [[ $environment = 'test' ]]; then
    repository="cr.yandex/crppdu233q7oqs81a1ta/"
elif [[ $environment = 'prod' ]]; then 
    repository="cr.yandex/crp7eamd5mno5u6gg0p9/"
fi

images=(
    "${repository}identity:$tag" 
    "${repository}file-storage:$tag" 
    "${repository}messaging:$tag" 
    "${repository}messaging-pg-migrate:$tag" 
    "${repository}user:$tag"
    "${repository}user-pg-migrate:$tag"
    "${repository}sms-service-stub:$tag"
)
paths=(
    "identity-service/" 
    "file-storage-service/" 
    "messaging-service/" 
    "messaging-service/" 
    "user-service/"
    "user-service/"
    "stubs/sms-service-stub/"
)
dockerfiles=(
    "Dockerfile"
    "Dockerfile"
    "Dockerfile"
    "migrate.Dockerfile"
    "Dockerfile"
    "migrate.Dockerfile"
    "Dockerfile"
)

build-parallel-local() {
    echo "Bulding docker images parallely..."

    for i in "${!images[@]}"; do
        minikube image build "${paths[$i]}" \
            --build-arg GOARCH="$(uname -m)" \
            -f "${dockerfiles[$i]}" \
            -t "${images[$i]}" \
            &> /dev/null &
    done

    if ! wait; then
        echo "At least one build failed! Run it with --debug flag"
        exit 1
    fi

    echo "Successfully built docker images!"
}

build-debug-local() {
    echo "Bulding docker images..."
    for i in "${!images[@]}"; do
        if ! minikube image build "${paths[$i]}" \
            --build-arg GOARCH="$(uname -m)" \
            -f "${dockerfiles[$i]}" \
            -t "${images[$i]}"
        then
            echo "Building ${images[$i]} failed"
            exit 1
        fi
    done

    echo "Successfully built docker images!"
}

build-parallel() {
    for i in "${!images[@]}"; do
        docker build "${paths[$i]}" \
            --build-arg GOARCH=amd64 \
            -f "${dockerfiles[$i]}" \
            -t "${images[$i]}" \
            &> /dev/null &
    done

    if ! wait; then
        echo "At least one build failed! Run it with --debug flag"
        exit 1
    fi

    echo "Successfully build docker images!"
    if $nopush ; then 
        return
    fi

    echo "Pushing images to the remote repository..."

    for i in "${!images[@]}"; do
        docker push "${images[$i]}" &> /dev/null &
    done

    if ! wait; then
        echo "At least one push failed! Run it with --debug flag"
        exit 1
    fi

    echo "Successfully pushed docker images..."
}

build-debug() {
    echo "Bulding docker images..."
    for i in "${!images[@]}"; do
        if ! docker build "${paths[$i]}" \
            --build-arg GOARCH=amd64 \
            -f "${dockerfiles[$i]}" \
            -t "${images[$i]}"
        then
            echo "Building ${images[$i]} failed"
            exit 1
        fi
    done

    echo "Successfully built docker images!"

    if $nopush ; then 
        return
    fi

    echo "Pushing images to the remote repository..."

    for i in "${!images[@]}"; do
        if ! docker push "${images[$i]}"; then
            echo "Pushing ${images[$i]} failed"
            exit 1
        fi
    done

    echo "Successfully pushed docker images..."
}

if [[ $environment == 'local' ]]; then
    if $debug; then 
        build-debug-local
    else
        build-parallel-local
    fi
else 
    if $debug; then 
        build-debug
    else
        build-parallel
    fi
fi
