# Solution

## Todo:

* general idea:
  1. solution with all parts and configs
  2. break topics appart 
  3. build chapters from the beginning towards the solution
  3.1. add introduction chapters with simplified examples if necessary

* kubernetes configs:
  * scheduler job (run to complete)

* gitlab-ci for build and deployment (on master)
  * use google cluster
  * use gitlab registry

* code documentation golang

* scaling example of crawler - as soon as crawler running, start next crawler

* check cluster login auf registry (gitlab registry

* write tests golang: xyz_test.go

* check commands on windows (\ problem new line)

* linux virtual machine fully prepared

* nice to have / check
  * migrations (change database field example)
  * ssl for backend

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
# this walthrough expects you to have everything installed (see INSTALL.md)

cd solution/code/

# cleanup if already running
docker network rm letsboot

for container in letsboot-backend letsboot-queue letsboot-database letsboot-frontend letsboot-crawler; do \
  docker stop $container
  docker container rm $container
done

kubectl get deployments --namespace letsboot
for deployment in frontend backend crawler; do \
  kubectl delete deployment letsboot-$deployment --namespace letsboot
done

kubectl get statefulset --namespace letsboot
helm delete letsboot-database -n letsboot
helm delete letsboot-queue -n letsboot

kubectl get pvc --namespace letsboot
for volume in $(kubectl get pvc --namespace letsboot -o name); do \
  kubectl delete $volume --namespace letsboot
done

# create a docker network for the containers to talk in
docker network create letsboot

# rabbitmq 
# note: the hostname in this case is only for rabbitmq important, for networking we use the --name
docker run -d --hostname rabbitmq --name letsboot-queue \
  -p 5672:5672 --network letsboot rabbitmq:3 \
  -e RABBITMQ_DEFAULT_PASS="megasecure" \
  -e RABBITMQ_DEFAULT_USER=letsboot

# mariadb - directly creates database and user
docker run --name letsboot-database \
  -e POSTGRES_PASSWORD="supersecure" \
  -e POSTGRES_USER="letsboot" \
  -e POSTGRES_DB="letsboot" \
  -p 3306:3306 -d --network letsboot postgres

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

# hint: if you build the images again, you'll see how much faster it is due to caching

# build frontend
docker build -t letsboot-frontend -f build/package/frontend.Dockerfile .

# run frontend
docker run -d --name letsboot-frontend --network letsboot \
  -p 4201:80  letsboot-frontend 

# open frontend http://localhost:4201/

# show your local images
docker images|grep letsboot

# show runing services
docker ps

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
mysql -p # enter password "supersecure!!"
SHOW DATABASES; # show databases

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

# run with database access
./backend --db.password=letsboot --queue.password=guest &
./crawler --queue.password=guest &

# push docker images to gitlab registry
# create token with registry_read und registry_write https://gitlab.com/profile/personal_access_tokens
# use email address and token to login to registry:

docker login registry.gitlab.com

# docker build -t registry.gitlab.com/letsboot/core/kubernetes-course .
# docker push registry.gitlab.com/letsboot/core/kubernetes-course

docker tag letsboot-backend registry.gitlab.com/letsboot/core/kubernetes-course/backend
docker tag letsboot-crawler registry.gitlab.com/letsboot/core/kubernetes-course/crawler
docker tag letsboot-frontend registry.gitlab.com/letsboot/core/kubernetes-course/frontend

# warning if you do not specify a private registry docker 
# may push your image to the public registry

docker push registry.gitlab.com/letsboot/core/kubernetes-course/backend
docker push registry.gitlab.com/letsboot/core/kubernetes-course/crawler
docker push registry.gitlab.com/letsboot/core/kubernetes-course/frontend

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

docker push eu.gcr.io/letsboot/kubernetes-course/backend
docker push eu.gcr.io/letsboot/kubernetes-course/crawler
docker push eu.gcr.io/letsboot/kubernetes-course/frontend

# run everything on kubernetes

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

# hint: we now use the passwords directly from the secrets
#       which are set by the helm statefullsets

# example: show secrets

$(kubectl get secret --namespace letsboot letsboot-queue -o jsonpath="{.data.rabbitmq-password}" | base64 --decode)
$(kubectl get secret --namespace letsboot letsboot-database-mariadb -o jsonpath="{.data.mariadb-password}" | base64 --decode)

# demo: how to run busybox on kubernetes
# this is a great way to log into a shell inside your kubernetes namespace
kubectl run -i --tty busybox --image=busybox --restart=Never --namespace letsboot -- sh

# list current pods (like containers but kubernetes :-)
kubectl get pods --namespace letsboot 

# manually deploy backend with it's service
kubectl apply -f deployments/backend/deployment.yaml --namespace letsboot
kubectl apply -f deployments/backend/service.yaml --namespace letsboot

# show logs of backend deployment
kubectl logs --selector=app=backend --namespace letsboot

kubectl apply -f deployments/frontend/deployment.yaml --namespace letsboot
kubectl apply -f deployments/frontend/service.yaml --namespace letsboot

kubectl apply -f deployments/crawler/deployment.yaml --namespace letsboot

# start deployment with backend, frontend and crawler at once using kustomize
# hint: kustomize is standard in kubernetes to have adaptable deployment configurations
kubectl apply -k deployments

# hint: to delete deployments
# kubectl delete deployment backend
# kubectl delete service backend

# per default networking is possible only inside cluster
# to access your services from outside you either have to configure a so called ingress
# or you can use port forwarding which we use untill we have ingress or if we
# want to access a service which doesn't need external access like postgres, rabbitmq...

# let's port forward the backend to use within the frontend
kubectl port-forward --namespace letsboot backend-b5c4fb56-5q2sh 8080:8080

# let's expse the frontend
kubectl port-forward --namespace letsboot frontend-856f54ddb4-9cbpk 4201:80


```