#!/bin/bash

# Building images via minikube to include it to its registry automatically
echo "Bulding docker images..."
minikube image build identity-service/ -t identity:local  
minikube image build messaging-service/ -t messaging:local  
minikube image build file-storage-service/ -t file-storage:local
minikube image build user-service/ -t user:local