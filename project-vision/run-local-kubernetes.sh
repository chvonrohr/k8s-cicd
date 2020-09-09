#!/bin/bash

./cleanup.sh

kubectl config use-context docker-desktop

parallel --jobs 4 << EOF
docker build -t letsboot-backend -f build/package/backend.Dockerfile .
docker build -t letsboot-crawler -f build/package/crawler.Dockerfile .
docker build -t letsboot-scheduler -f build/package/scheduler.Dockerfile .
docker build -t letsboot-frontend -f build/package/frontend.Dockerfile .
EOF

docker tag letsboot-backend eu.gcr.io/letsboot/kubernetes-course/backend
docker tag letsboot-crawler eu.gcr.io/letsboot/kubernetes-course/crawler
docker tag letsboot-frontend eu.gcr.io/letsboot/kubernetes-course/frontend
docker tag letsboot-scheduler eu.gcr.io/letsboot/kubernetes-course/scheduler

docker push eu.gcr.io/letsboot/kubernetes-course/backend
docker push eu.gcr.io/letsboot/kubernetes-course/crawler
docker push eu.gcr.io/letsboot/kubernetes-course/frontend
docker push eu.gcr.io/letsboot/kubernetes-course/scheduler

kubectl create namespace letsboot

kubectl config set-context --current --namespace=letsboot

helm repo add bitnami https://charts.bitnami.com/bitnami

helm install letsboot-queue --set replicaCount=3 bitnami/rabbitmq -n letsboot 
helm install letsboot-database --set global.postgresql.postgresqlDatabase=letsboot,global.postgresql.postgresqlUsername=letsboot bitnami/postgresql -n letsboot

kubectl apply -k deployments

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/cloud/deploy.yaml

kubectl apply -f deployments/web-ingress.yaml

open http://localhost