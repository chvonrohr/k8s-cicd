# our project in docker

* backend
* frontend
* crawler
* scheduler
* database
* messagequeue

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
  --name letsboot-queue \
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
  --name letsboot-database \
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
ng serve -o
ng build --prod
ls dist/crawler/
```

----

## Frontend - docker walkthrough

1. use node image
2. copy dependencie info
3. install dependencies
4. copy code
5. build app
6. use nginx image
7. copy build
8. serve

----

## Frontend - multistage 1/2

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

project/
```bash
docker build -t letsboot-frontend \
  -f build/package/frontend.Dockerfile .

docker run -d --name letsboot-frontend --network letsboot \
  -p 4201:80  letsboot-frontend 
```

http://localhost:4201/

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

## Backedn - manual

```bash
go mod download
go build ./cmd/backend
./backend --db.password="supersecure" --queue.password="megasecure" 
```

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
* As golang is statically compiled we don't need anything

----

## Test app 

1. https://localhost:4201
2. add website

----

### crawler


---- 

### scheduler

* Will be used in kubernetes.


---

### shutdown everything



----

### gitlab registry
