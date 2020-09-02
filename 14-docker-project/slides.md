# Our Project in Docker

* backend - core logic and api
* frontend - add websites and see data
* crawler - listens to the queue and crawls website
* scheduler - triggers crawler jobs
* database - store data
* messagequeue - manage jobs in queue

----

## our plan

1. create network
2. run RabbitMQ
3. run PostgreSQL
4. build containers
5. run containers
6. ? deploy to registry
7. ? build in gitlab-ci

----

## create network

```bash
# create a docker network for the containers to talk in
docker network create letsboot
```

----

## RabbitMQ

```bash
# rabbitmq 
docker run -d \
  --name queue \
  --hostname rabbitmq \
  --network letsboot \
  -p 5672:5672 \
  -e RABBITMQ_DEFAULT_PASS="megasecure" \
  -e RABBITMQ_DEFAULT_USER=letsboot \
  rabbitmq
```

Note: 
* The hostname in this case is only for rabbitmq important, for networking we use the --name

----

## PostgreSQL

```bash
# postgresql - directly creates database and user
docker run -d \
  --name database \
  --network letsboot \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD="supersecure" \
  -e POSTGRES_USER="letsboot" \
  -e POSTGRES_DB="letsboot" \
  postgres
```

----

## Frontend - manual build

```bash
cd web
yarn install
ng serve --host 0.0.0.0 --disable-host-check # ctrl+c to exit
echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4200/
ng build --prod
ls dist/crawler/
```

----

## Frontend - docker walkthrough

1. use node image
2. copy dependency info
3. install dependencies
4. copy code
5. build app
6. use nginx image
7. copy build to nginx

----

## Frontend - container multistage 1/2

web/.dockerignore
```txt
node_modules
```

build/ci/package/frontend.Dockerfile
```Dockerfile
FROM node:12 as build
WORKDIR /app
COPY web/yarn.lock .
COPY web/package.json .
RUN yarn install
COPY web/ .
RUN node_modules/.bin/ng build --prod
# ...
```

----

## Frontend - multistage 2/2

build/ci/package/frontend.Dockerfile
```Dockerfile
FROM nginx:stable
COPY --from=build /app/dist/crawler /usr/share/nginx/html
```

----

## Frontend - build and run

project-start/
```bash
docker build -t frontend \
  -f build/package/frontend.Dockerfile .

docker run -d --name letsboot-frontend \
  --network letsboot -p 4201:80 frontend 

echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4201/
```

----

## Backend - multistage 1/2 

build/ci/package/backend.Dockerfile
```Dockerfile
FROM golang:alpine AS build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build ./cmd/backend
# ...
```

Note:
* we will use a scratch image, as go is statically compiled and doesn't need anything


----

## Backend - multistage 2/2 scratch

`scratch` the empty base image of everything

build/ci/package/backend.Dockerfile
```Dockerfile
# ...
FROM scratch
WORKDIR /app 
COPY --from=build /app/backend /app/backend
COPY --from=build /app/config/backend.* .
ENTRYPOINT ["/app/backend"]
```

Note: 
* As golang is statically compiled we don't need anything else

----

## Backend - build and run

project-start/
```bash
docker build -t backend \
  -f build/package/backend.Dockerfile .

docker run -d --name backend -p 8080:8080 \
  -e LETSBOOT_DB.HOST=database \
  -e LETSBOOT_DB.PASSWORD="supersecure" \
  -e LETSBOOT_QUEUE.HOST=queue \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  -v /home/letsboot/docker-volume/:/var/data \
  --network letsboot backend

echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:8080/sites
```

----

## Test app 

```bash
echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4201/ 
echo add website
echo check data: http://$PARTICIPANT_NAME.sk.letsboot.com:8080/sites
```

----

### Crawler - Dockerfile

project-start/
```bash
cp build/package/backend.Dockerfile build/package/crawler.Dockerfile
```

build/package/crawler.Dockerfile
```Dockerfile
#...
RUN CGO_ENABLED=0 go build ./cmd/crawler
FROM scratch
WORKDIR /app 
COPY --from=build /app/crawler /app/crawler
COPY --from=build /app/config/backend.* .
ENTRYPOINT ["/app/crawler"]
#...
```
----

## Crawler - build and run

project-start/
```bash
docker build -t crawler \
  -f build/package/crawler.Dockerfile .

docker run -d --name crawler  \
  -e LETSBOOT_QUEUE.HOST=queue \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  -v /home/letsboot/docker-volume/:/var/data \
  --network letsboot crawler
```

> no port as the crawler listens to the queue

----

### Crawler - manually invoke crawling

```bash
curl -H "Content-Type: application/json" \
    -X POST -d '{"siteId":1}' \
    http://localhost:8080/crawls
```

---- 

### scheduler

build/package/scheduler.Dockerfile
```Dockerfile
FROM curlimages/curl
CMD curl $SCHEDULE_URL
```

project-start/
```bash
# build without context
docker build -t scheduler - < build/package/scheduler.Dockerfile 

# this will be run with kubernetes jobs
docker run -it -e SCHEDULE_URL=http://backend:8080/schedule --network letsboot scheduler
```

Note:
* try `http://backend:8080/sites` to see a result

---

### Shutdown everything

```sh
docker stop frontend backend crawler database queue
docker rm frontend backend crawler database queue
```

----

### gitlab registry

```bash
docker tag backend registry.gitlab.com/$GIT_REPO/jonasfelix/backend:latest
docker tag crawler registry.gitlab.com/$GIT_REPO/jonasfelix/crawler:latest
docker tag frontend registry.gitlab.com/$GIT_REPO/jonasfelix/frontend:latest
docker tag scheduler registry.gitlab.com/$GIT_REPO/jonasfelix/scheduler:latest

docker push registry.gitlab.com/$GIT_REPO/jonasfelix/backend:latest
docker push registry.gitlab.com/$GIT_REPO/jonasfelix/crawler:latest
docker push registry.gitlab.com/$GIT_REPO/jonasfelix/frontend:latest
docker push registry.gitlab.com/$GIT_REPO/jonasfelix/scheduler:latest

echo "open https://gitlab.com/$GIT_REPO/container_registry"
```

----

### recap

* run prebuilt docker images (database, queue)
* create Dockerfiles
* run containers
* stop remove containers
* push containers to registry