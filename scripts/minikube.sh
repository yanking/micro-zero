#!/bin/bash

echo "=============> Enabling registry addon"
minikube addons enable registry

echo "=============> Running socat"
docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"
