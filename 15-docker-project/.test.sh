#!/bin/bash

startdirectory=$(pwd)

if [ ! -d "./project-start" ]; then \
    echo "run in course folder" 1>&2
    exit 1
fi 

function cleanuptest() {
    cd $startdirectory
    docker rm -f database queue frontend backend crawler
    docker network rm letsboot
    docker volume rm page-storage
    rm -r project-start
    mv project-start-pretest project-start
}

function errout() {
    echo Error: "$1" 1>&2
    cleanuptest
    exit 1
}

cp -r project-start project-start-pretest
cp -rv 15-docker-project/solution/* project-start/ ||errout "merging solution"

cd project-start ||errout "no project start folder"

parallel -j4 << EOF
docker build -t crawler -f build/package/crawler.Dockerfile .
docker build -t scheduler -f build/package/scheduler.Dockerfile .
docker build -t backend -f build/package/backend.Dockerfile .
docker build -t frontend -f build/package/frontend.Dockerfile .
EOF

docker network create letsboot ||errout "creating network"

docker run -d \
  --name database \
  --network letsboot \
  -e POSTGRES_PASSWORD="supersecure" \
  -e POSTGRES_USER="letsboot" \
  -e POSTGRES_DB="letsboot" \
  postgres ||errout "running database"

docker run -d \
  --name queue \
  --hostname rabbitmq \
  --network letsboot \
  -e RABBITMQ_DEFAULT_PASS="megasecure" \
  -e RABBITMQ_DEFAULT_USER=letsboot \
  rabbitmq ||errout "running queue"

sleep 10

docker volume create page-storage||errout "couldn't create storage"

docker run -d --name frontend \
  --network letsboot -p 4201:80 frontend ||errout "running frontend"

docker run -d --name backend -p 8080:8080 \
  -e LETSBOOT_DB.HOST=database \
  -e LETSBOOT_DB.PASSWORD="supersecure" \
  -e LETSBOOT_DB.TYPE="postgres" \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  -e LETSBOOT_QUEUE.HOST=queue \
  -v page-storage:/var/data \
  --network letsboot backend ||errout "running backend"

docker run -d --name crawler  \
  -e LETSBOOT_BACKEND.URL=http://backend:8080 \
  -e LETSBOOT_QUEUE.HOST=queue \
  -e LETSBOOT_QUEUE.PASSWORD="megasecure" \
  -v page-storage:/var/data \
  --network letsboot crawler  ||errout "running crawler"

sleep 10

docker run -e SCHEDULE_URL=http://backend:8080/schedule \
  --network letsboot scheduler  ||errout "running scheduler"

curl http://localhost:8080 ||errout "no backend"
curl http://localhost:4201 ||errout "no frontend"

docker tag backend registry.gitlab.com/$GIT_REPO/backend:latest ||errout "tag"
docker tag crawler registry.gitlab.com/$GIT_REPO/crawler:latest ||errout "tag"
docker tag frontend registry.gitlab.com/$GIT_REPO/frontend:latest ||errout "tag"
docker tag scheduler registry.gitlab.com/$GIT_REPO/scheduler:latest ||errout "tag"

docker push registry.gitlab.com/$GIT_REPO/backend:latest ||errout "push"
docker push registry.gitlab.com/$GIT_REPO/crawler:latest ||errout "push"
docker push registry.gitlab.com/$GIT_REPO/frontend:latest ||errout "push"
docker push registry.gitlab.com/$GIT_REPO/scheduler:latest ||errout "push"

cleanuptest