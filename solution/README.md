# Solution

## Todo:

* general idea:
  1. solution with all parts and configs
  2. break topics appart 
  3. build chapters from the beginning towards the solution
  3.1. add introduction chapters with simplified examples if necessary
  
## Application overview

* based on golang recommended structure:
  * https://github.com/golang-standards/project-layout

* parts
  * backend - rest api and core business logic
    * model: site, page, crawl
  * frontend - simple list of domains and websites
  * scheduler - simple script calling backend to schedule crawl jobs
  * crawler - service to crawl websites listening to queue
  * database - postgres
  * rabbitmq - message queue service

* code/
  * build/package/ - dockerfiles
  * cmd/ - golang entry points for services
  * deployments/ - kubernetes configuration   
  * internal/ - source code for backend and crawler
  * web/ - frontend
  * config/ - config files for services
    * backend.toml - config file for backend
    * crawler.toml - config file for crawler
    * default.conf - config for nginx server
  * go.mod - golang dependencies
  * go.sum - golang version locks
  * .dockerignore - files to ignore with COPY . . to prevent rebuilds (caching)

## notes

setup frontend

```
ng new crawler --prefix crl --style scss --skip-git --directory web --create-application false
cd web
ng generate application crawler --prefix crl --routing true --style scss
```

## walkthrough

```bash
# path cleanup.sh
```

```bash
# walkthrough start - do not remove -

# this walthrough expects you to have everything installed (see INSTALL.md)

cd solution/code/

# cleanup if already running
# alternatively you can run cleanup.sh

docker network rm letsboot

for container in $(docker ps --filter network=letsboot --format '{{.Names}}'); do \
  docker stop $container
  docker container rm $container
done

for deployment in $(kubectl get deployments --namespace letsboot -o name); do \
  kubectl delete $deployment --namespace letsboot
done

kubectl get statefulset --namespace letsboot
helm delete letsboot-database -n letsboot
helm delete letsboot-queue -n letsboot

for volume in $(kubectl get pvc --namespace letsboot -o name); do \
  kubectl delete $volume --namespace letsboot
done

# create a docker network for the containers to talk in
docker network create letsboot

# rabbitmq 
# note: the hostname in this case is only for rabbitmq important, for networking we use the --name
docker run -d --hostname rabbitmq --name letsboot-queue \
  -p 5672:5672 --network letsboot \
  -e RABBITMQ_DEFAULT_PASS="megasecure" \
  -e RABBITMQ_DEFAULT_USER=letsboot \
  rabbitmq:3 
  

# mariadb - directly creates database and user
docker run --name letsboot-database \
  -e POSTGRES_PASSWORD="supersecure" \
  -e POSTGRES_USER="letsboot" \
  -e POSTGRES_DB="letsboot" \
  -p 5432:5432 -d --network letsboot postgres

# build backend
docker build -t letsboot-backend -f build/package/backend.Dockerfile .

# run backend
docker run -d --name letsboot-backend -p 8080:8080 \
  -e LETSBOOT_DB.HOST=letsboot-database \
  -e LETSBOOT_QUEUE.HOST=letsboot-queue \
  -e LETSBOOT_DB.PASSWORD="supersecure" \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  --network letsboot letsboot-backend

# build crawler
docker build -t letsboot-crawler -f build/package/crawler.Dockerfile .

# run crawler
docker run -d --name letsboot-crawler \
  -e LETSBOOT_QUEUE.HOST=letsboot-queue \
  -e LETSBOOT_BACKEND.URL="http://letsboot-backend:8080" \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  --network letsboot letsboot-crawler

# build scheduler
docker build -t letsboot-scheduler -f build/package/scheduler.Dockerfile .

# run scheduler (once) - with kubernetes we'll use the kubernetes scheduler
# we don't use -d "detached" here, as we want to directly see how it runs
docker run --name letsboot-scheduler --network letsboot \
  letsboot-scheduler "http://letsboot-backend:8080/schedule"

# hint: if you build the images again, you'll see how much faster it is due to caching

# build frontend
docker build -t letsboot-frontend -f build/package/frontend.Dockerfile .

# run frontend
docker run -d --name letsboot-frontend --network letsboot \
  -p 4201:80  letsboot-frontend 

# show your local images
docker images|grep letsboot

# show runing services
docker ps

# open frontend in browser http://localhost:4201/
open http://localhost:4201/

# test backend service
# note: ssl will be done by the revese proxy in kubernetes
curl http://localhost:8080/sites

# add site on commandline
curl -H "Content-Type: application/json" \
   -X POST -d '{"url":"https://www.letsboot.com","interval":3600000}' \
   http://localhost:8080/sites

# check added site - this way you see that database works
curl http://localhost:8080/sites

# in a separat terminal window stream logs of crawler
docker logs -f letsboot-crawler

# start crawler - this adds an item to the rabbitmq and the crawler will pick it up
curl -H "Content-Type: application/json" \
   -X POST -d '{"siteId":1}' \
   http://localhost:8080/crawls

# run to complete for https://localhost:8080/sites/crawl
# minimal busybox setup
# ... todo

# how to get a shell in the mariadb docker process
docker exec -it letsboot-database /bin/bash
# psql -U letsboot -W # enter password "supersecure!!"
# \list # show databases
# \dt # show tables

# how to start an interactive busybox container
docker run -it --network letsboot  busybox

# how to run backend and crawler localy (ie. for debugging with breakpoints)
docker stop letsboot-backend
docker stop letsboot-crawler

# get dependencies to build on host
go mod download

# backend (host)
go build ./cmd/backend

# crawler (host)
go build ./cmd/crawler

# crawler (host)
go build ./cmd/scheduler

# run with database access
./backend --db.password="supersecure" --queue.password="megasecure" &
./crawler --queue.password="megasecure" &
./scheduler "https://localhost:8080" 

# push docker images to gitlab registry
# create token with registry_read und registry_write https://gitlab.com/profile/personal_access_tokens
# use email address and token to login to registry:

docker login registry.gitlab.com

# docker build -t registry.gitlab.com/letsboot/core/kubernetes-course .
# docker push registry.gitlab.com/letsboot/core/kubernetes-course

docker tag letsboot-backend registry.gitlab.com/letsboot/core/kubernetes-course/backend
docker tag letsboot-crawler registry.gitlab.com/letsboot/core/kubernetes-course/crawler
docker tag letsboot-frontend registry.gitlab.com/letsboot/core/kubernetes-course/frontend
docker tag letsboot-scheduler registry.gitlab.com/letsboot/core/kubernetes-course/scheduler

# warning if you do not specify a private registry docker 
# may push your image to the public registry

docker push registry.gitlab.com/letsboot/core/kubernetes-course/backend
docker push registry.gitlab.com/letsboot/core/kubernetes-course/crawler
docker push registry.gitlab.com/letsboot/core/kubernetes-course/frontend
docker push registry.gitlab.com/letsboot/core/kubernetes-course/scheduler

# hint: as some layers are the same (like the first steps COPY ... in scratch) 
# not all layers have to be pushed three times, docker is extremly optimized in this point

# For fallback reasons we have a public google registry where only we can push
# this registry is used if you have troubles with your personal gitlab registry

gcloud auth
gcloud config set project letsboot
gcloud auth configure-docker eu.gcr.io

docker tag letsboot-backend eu.gcr.io/letsboot/kubernetes-course/backend
docker tag letsboot-crawler eu.gcr.io/letsboot/kubernetes-course/crawler
docker tag letsboot-frontend eu.gcr.io/letsboot/kubernetes-course/frontend
docker tag letsboot-scheduler eu.gcr.io/letsboot/kubernetes-course/scheduler

docker push eu.gcr.io/letsboot/kubernetes-course/backend
docker push eu.gcr.io/letsboot/kubernetes-course/crawler
docker push eu.gcr.io/letsboot/kubernetes-course/frontend
docker push eu.gcr.io/letsboot/kubernetes-course/scheduler

## --- run everything on kubernetes

# make sure you are in the correct context
# either through the menu of your local docker desktop
# or the following commands

# show contexts and selected one with *
kubectl config get-contexts

# get current context name
kubectl config current-context

# set local docker desktop context if not already the case
kubectl config use-context docker-desktop

# install dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended.yaml

# make dashobard available
kubectl proxy

# get a token
dashboard_token=$(kubectl -n kube-system describe secret default |grep "token:"|awk '{ print $2 }') 
kubectl config set-credentials docker-for-desktop --token="$dashboard_token"
echo $dashboard_token

# open in browser and copy the token to authenticate the dashboard
open "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login"


# create namespace for kubernetes
# note: only the images are used from docker, everything else is separate
kubectl create namespace letsboot

# set letsboot as our current namespace
kubectl config set-context --current --namespace=letsboot

# hint: the rabbbitmq and postgres setups we use on kubernetes are NOT the
# same as on docker, as we want clustering and management of statefull 
# sets which we don't have in docker

# add bitnami for our rabbitmq and postgres setups
helm repo add bitnami https://charts.bitnami.com/bitnami

# get statefullset of rabbitmq and run it
helm install letsboot-queue --set replicaCount=3 bitnami/rabbitmq -n letsboot 

# get statefullset of postgres and run it
helm install letsboot-database --set global.postgresql.postgresqlDatabase=letsboot,global.postgresql.postgresqlUsername=letsboot bitnami/postgresql -n letsboot

# more about scaling and replicas of postgres here: https://github.com/bitnami/charts/tree/master/bitnami/postgresql

# hint: we now use the passwords directly from the secrets
#       which are set by the helm statefullsets

# example: show secrets

kubectl get secret --namespace letsboot letsboot-queue-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 --decode
kubectl get secret --namespace letsboot letsboot-database-postgresql -o jsonpath="{.data.postgresql-password}" | base64 --decode

# demo: how to run busybox on kubernetes
# this is a great way to log into a shell inside your kubernetes namespace
kubectl run -i --tty busybox --image=busybox --restart=Never --namespace letsboot -- sh

# list current pods (like containers but kubernetes :-)
# you should see the database and the queue
kubectl get pods --namespace letsboot 

# example how to apply specific configurations
# kubectl apply -f deployments/frontend/deployment.yaml --namespace letsboot
# kubectl apply -f deployments/frontend/service.yaml --namespace letsboot

# start deployment with backend, frontend and crawler at once using kustomize
# hint: kustomize is standard in kubernetes to have adaptable deployment configurations
kubectl apply -k deployments

# show logs of backend deployment
kubectl logs --selector=app=backend --namespace letsboot

# hint: to delete deployments
# kubectl delete deployment backend
# kubectl delete service backend

# per default networking is possible only inside cluster
# to access your services from outside you either have to configure a so called ingress
# or you can use port forwarding which we use untill we have ingress or if we
# want to access a service which doesnt need external access like postgres, rabbitmq

# how to get the first pod of a app
# kubectl get pods -l app=backend -o name|head -n1

# let's port forward the backend to use within the frontend
kubectl port-forward --namespace letsboot service/letsboot-backend 8080:80

# let's expse the frontend
# trick question: why can we not only expose the frontend? why do we need to expose the backend?
kubectl port-forward --namespace letsboot service/letsboot-frontend 4201:80

# alternative let's take a sneak peak at ingress by installing a local ingress controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/cloud/deploy.yaml

# we prepareted a ingress configuration for localhost
kubectl apply -f deployments/local-ingress.yaml
kubectl get ingress


## -------- Google Cluster


# get google credentials for exising cluster - automatically done on creation
# gcloud container clusters get-credentials gke_letsboot_europe-west6_jonas1

# let's create our own cluster
gcloud container clusters create jonas1 --project letsboot --region europe-west6 --machine-type e2-small --num-nodes 1

# allow kubernetes cluster to use your gitlab registry
# for simplicty reasons we use the public google registry for the training to reduce the amount of authentication
# 
# kubectl create secret docker-registry regcred --docker-server=<your-registry-server> --docker-username=<your-name> --docker-password=<your-pword> --docker-email=<your-email>
# 
# to use gitlab registry change the deployment.yaml files of each pod to (image:)
#       containers:
#        - name: backend
#          image: eu.gcr.io/letsboot/kubernetes-course/backend:latest

# for more information about options
# gcloud container clusters create --help

# switch context to the new cluster
kubectl config use-context gke_letsboot_europe-west6_jonas1

# check if there are any pods (should be empty)
kubectl get pods -n letsboot

# the same as above
kubectl create namespace letsboot 
kubectl config set-context --current --namespace=letsboot

# deploy stateful sets using helm
helm install letsboot-queue --set replicaCount=3 bitnami/rabbitmq -n letsboot 
helm install letsboot-database --set global.postgresql.postgresqlDatabase=letsboot,global.postgresql.postgresqlUsername=letsboot bitnami/postgresql -n letsboot

# applay deployments
kubectl apply --kustomize deployments
kubectl get pods --namespace letsboot

# expose the services from your google cluster to your local system
kubectl port-forward --namespace letsboot service/letsboot-backend 8080:80 & 
kubectl port-forward --namespace letsboot service/letsboot-frontend 4201:80 &

# show environment variables you get as a pod
kubectl run -i --tty busybox --image=busybox --restart=Never --namespace letsboot -- env; kubectl delete pod busybox

# experiment 
# let's create some load
kubectl top pods

for i in {1..100}; do \
  curl -H "Content-Type: application/json" \
    -X POST -d '{"url":"https://www.letsboot.com","interval":3600000}' \
    http://localhost:8080/sites
  curl http://localhost:8080/sites
done

# one of the backends and the database will slightly increase in usage
# 5mb for a backend ;-)
kubectl top pods

# let's start crawling the first 100 websites
for i in {1..100}; do \
  curl -H "Content-Type: application/json" \
    -X POST -d "{\"siteId\":$i}" \
    http://localhost:8080/crawls
done

# see difference (1m = 0.1% of a vcpu)
kubectl top pods

# show urls
curl http://localhost:8080/pages

# show logs of crawler
kubectl logs --selector=app=crawler --namespace letsboot

# install dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended.yaml

# make dashobard available
kubectl proxy

# get a token
dashboard_token=$(kubectl -n kube-system describe secret default |grep "token:"|awk '{ print $2 }') 
kubectl config set-credentials docker-for-desktop --token="$dashboard_token"
echo $dashboard_token

# open in browser and copy the token to authenticate the dashboard
open "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login"

# delete cluster
gcloud container clusters delete jonas1 --project letsboot --region europe-west6

# walkthrough end - do not remove -
```