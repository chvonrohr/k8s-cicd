# Solution

## Todo:

* write tests for everything: xyz_test.go
* 

## Files & Folder

* based on golang recommended structure:
  * https://github.com/golang-standards/project-layout

* code/
  * go.mod / go.sum - dependencies and version locks
  * build/package/ - dockerfiles
  * web/ - frontend
  * deployments/ - kubernetes configuration
  * internal/ - golang projects / code

## Setup solution

```
cd code/
go mod download

# use ./ in go build - otherwise it will try to build a package
go build ./cmd/crawler # creates crawler binary in current folder
go build ./cmd/backend # creates backend binary in current folder

# build docker image which is put into your local docker repository
docker build -t letsboot-backend -f build/package/backend.Dockerfile .

# show your local image
docker images|grep letsboot

# to see cahing run the build again and look at which parts are cached like the go mod download
```

## Docker build process

1. 

## Kubernetes Deployment

    * MariaDB - ?
    * RabbitMQ - simple as possilble
    * Golang REST Backend - ./backend
      * new site => add to sql write to rabbitmq with site id
      * db: 
        * website (id, starting_url, interval)
        * urls (id, website_id, url)
    * Angular Minimal-Frontend - ./frontend
      * input: "add website*
      * list with added websites
      * click on website shows list of urls
    * Golang Minimal Crawler  - ./crawler
      * listens to rabbitmq
      * crawls page 
      * sends found url to rest of backend
    * Kubernetes Cronjob mit Run to complete - ./scheduler
      * shellscript wget backend/schedule => ads pages by interval to rabbitmq

## Deploy locally


### Docker

```bash
# create a docker network for the containers to talk in
docker network create letsboot

# rabbitmq
docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 --network letsboot rabbitmq:3
# mariadb (make sure to create a "letsboot" database after creating this container)
docker run --name some-mariadb \
  -e MYSQL_ROOT_PASSWORD="test" \
  -e MYSQL_USER="letsboot" \
  -e MYSQL_PASSWORD="letsboot" \
  -e MYSQL_DATABASE="letsboot" \
  -p 3306:3306 -d --network letsboot mariadb

# backend (host)
go build ./cmd/backend && ./backend
# crawler (host)
go build ./cmd/crawler && ./crawler

# don't forget to amend configuration files for this

# backend (docker)
docker build -t letsboot-backend -f build/package/backend.Dockerfile .
docker run -d --name letsboot-backend -p 8080:8080 \
  -e LETSBOOT_DB.HOST=some-mariadb \
  -e LETSBOOT_QUEUE.HOST=some-rabbit \
  --network letsboot letsboot-backend

# crawler (docker)
docker build -t letsboot-crawler -f build/package/crawler.Dockerfile .
docker run -d --name letsboot-crawler \
  -e LETSBOOT_QUEUE.HOST=some-rabbit \
  -e LETSBOOT_BACKEND.URL="http://letsboot-backend:8080" \
  --network letsboot letsboot-crawler
```