# Solution

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
# rabbitmq
docker run -d --hostname my-rabbit --name some-rabbit -p 5672:5672 rabbitmq:3
# mariadb (change password)
docker run --name some-mariadb -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mariadb

# backend (host)
go build ./cmd/backend && ./backend
# crawler (host)
go build ./cmd/crawler && ./crawler

# don't forget to amend configuration files for this
# backend (docker)
docker build -t letsboot-backend -f build/package/backend.Dockerfile .
docker run -d --name letsboot-backend -p 8080:8080 letsboot-backend
# crawler (docker)
docker build -t letsboot-crawler -f build/package/crawler.Dockerfile .
docker run -d --name letsboot-crawler letsboot-crawler
```