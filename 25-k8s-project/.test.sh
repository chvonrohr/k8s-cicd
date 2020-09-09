#!/bin/bash

startdirectory=$(pwd)

if [ ! -d "./project-start" ]; then \
    echo "run in course folder" 1>&2
    exit 1
fi 

function cleanuptest() {
    cd $startdirectory
    kubectl delete namespace letsboot-test 
    kubectl config set-context --current --namespace=default
    rm -r project-start
    mv project-start-pretest project-start
}

function errout() {
    echo Error: "$1" 1>&2
    cleanuptest
    exit 1
}

kubectl config set-context kind-kind ||errout "switch context"

cp -r project-start project-start-pretest ||errout "backup start"
cp -rv 25-k8s-project/solution/* project-start/ ||errout "merging solution"

kubectl create namespace letsboot-test ||errout "create namespace"
kubectl config set-context --current --namespace=letsboot-test ||errout "set namespace"

kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=ItsComplicated! ||errout "database secret"

kubectl create secret generic queue-rabbitmq \
  --from-literal=rabbitmq-password=MoreSecrets! ||errout "queue secret"

kubectl apply -Rf project-start/deployments
kubectl apply -Rf project-start/kind-ingress.yaml

for deployment in queue database backend frontend crawler; do \
    kubectl wait deployments/$deployment --for condition=available ||errout "deployment $deployment not available"
done

curl http://localhost/api/ ||errout "no backend"
[ "$(curl http://localhost/api)" == "backend works" ] || ||errout "backend doesn't indicate to work"
curl http://localhost ||errout "no frontend"

cleanuptest