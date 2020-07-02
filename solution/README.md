# Solution

## Todo:

* write tests for everything: xyz_test.go
* code documentation
* scaling with example of crawler
* nice to have / check
  * migrations (change database field example)
  * ssl for backend

## Files & Folder

* based on golang recommended structure:
  * https://github.com/golang-standards/project-layout

* code/
  * build/package/ - dockerfiles
  * cmd/ - golang entry points for services
  * deployments/ - kubernetes configuration
  * internal/ - source code for backend and crawler
  * web/ - frontend
  * backend.toml - config file for backend
  * crawler.toml - config file for crawler
  * go.mod - golang dependencies
  * go.sum - golang version locks
  * .dockerignore - files to ignore with COPY . . to prevent rebuilds (caching)


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


## Setup and build locally without docker

### Docker

```bash
cd solution/code/

# create a docker network for the containers to talk in
docker network create letsboot

# rabbitmq 
# note: the hostname in this case is only for rabbitmq important, for networking we use the --name
docker run -d --hostname rabbitmq --name letsboot-queue -p 5672:5672 --network letsboot rabbitmq:3

# mariad
docker run --name letsboot-database \
  -e MYSQL_ROOT_PASSWORD="test" \
  -e MYSQL_USER="letsboot" \
  -e MYSQL_PASSWORD="letsboot" \
  -e MYSQL_DATABASE="letsboot" \ # directly creates a database
  -p 3306:3306 -d --network letsboot mariadb

# todo: commands to run without docker

# don't forget to amend configuration files for this

# build backend
docker build -t letsboot-backend -f build/package/backend.Dockerfile .

# run backend
docker run -d --name letsboot-backend -p 8080:8080 \
  -e LETSBOOT_DB.HOST=letsboot-database \
  -e LETSBOOT_QUEUE.HOST=letsboot-queue \
  -e LETSBOOT_DB.PASSWORD=letsboot \
  -e LETSBOOT_QUEUE.PASSWORD=guest \
  --network letsboot letsboot-backend

# build crawler
docker build -t letsboot-crawler -f build/package/crawler.Dockerfile .

# run crawler
docker run -d --name letsboot-crawler \
  -e LETSBOOT_QUEUE.HOST=letsboot-queue \
  -e LETSBOOT_BACKEND.URL="http://letsboot-backend:8080" \
  -e LETSBOOT_QUEUE.PASSWORD=guest \
  --network letsboot letsboot-crawler

# hint: if you build the images again, you'll see how much faster it is due to caching

# to restart a docker image you have to do the following
docker stop letsboot-crawler
docker rm letsboot-crawler
# ... build and run as seen above

# show your local images
docker images|grep letsboot

# show runing services
docker ps

# test backend service
# note: ssl will be done by the revese proxy in kubernetes
http://localhost:8080/sites

# add site on commandline
curl -H "Content-Type: application/json" \
   -X POST -d '{"url":"https://www.letsboot.com","interval":3600000}' \
   http://localhost:8080/sites

# check added site - this way you see that database works
curl http://localhost:8080/sites

# in a separat terminal window stream logs of crawler
docker logs -f letsboot-crawler

# start crawler - this adds an item to the rabbitmq and the crawler will pick it up
curl -X POST http://localhost:8080/sites/1/crawl

# run to complete for https://localhost:8080/sites/crawl
# minimal busybox setup
... todo

# how to start ainteractive busybox container
docker run -it --network letsboot  busybox

# how to get a shell in the mariadb docker process
docker exec -it letsboot-database /bin/bash
mysql -p # enter password
SHOW DATABASES; # show databases

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

```




## Deploy kubernetes locally