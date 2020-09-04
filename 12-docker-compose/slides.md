# Docker compose

* each container should do one thing (process)
* applications need more than one process

----

## let's add MariaDB to our app

```txt
[todo-app]--[MariaDB]
```

* create network
* start mariadb
* run app with mariadb

----

## create a network

* container networks are isolated
* container in the same network can connect to each other

```bash
docker network create todo-app
```

----

## start mariadb in network

```bash
# create volume
docker volume create todo-mariadb

# run a database
docker run -d \
  --network todo-app --network-alias database
  -v todo-mariadb:/var/lib/mysql
  -e MYSQL_ROOT_PASSWORD=supersecret \
  -e MYSQL_DATABASE=todos \
  mariadb:10

# list databases
docker exec -it 75cb mysql --password=supersecret -e "SHOW DATABASES;"
```
> in production we'll use secretes to provide access information

Note:
* `--network todo-app` use our new network "todo"
* `--network-alias database` let us connect to the mariadb using database as hostname
* `-e MYSQL_ROOT_PASSWORD=supersecret` set an environment variable for within the container

----

## look into the network

```bash
# use a great debugging image
docker run -it --network todo-app nicolaka/netshoot

# show ip of the database
dig database
```

Note:
Docker provides a internal DNS server for container lookups.

----

## run our minimal app

Configure database through env variables.
```bash
docker run -dp 4000:3000 \
  --network todo-app \
  -e MYSQL_HOST=database \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD=supersecret \
  -e MYSQL_DB=todos \
  todo-app
```

----

## see database changes

Open localhost:3000 and add some todos

```bash
docker exec -it 75cb mysql --password=supersecret \
    -e "USE todos; SHOW TABLES;"

docker exec -it 75cb mysql --password=supersecret \
    -e "USE todos; SELECT * FROM todo_items;"
```

----

## now let's restart everything

![complicated](https://media.giphy.com/media/cnuQwZ8IFLDZFwreWF/giphy.gif)

----

## docker compose file

todo-app/docker-compose.yml
```yaml
version: "3.7"

services:
  app:
    image: minimal-app
    ports:
      - 4000:3000
    environment:
      MYSQL_HOST: database
      MYSQL_USER: root
      MYSQL_PASSWORD: supersecret
      MYSQL_DB: todos
  database:
    image: mariadb:10
    volumes:
      - todo-mariadb:/var/lib/mysql
    environment: 
      MYSQL_ROOT_PASSWORD: supersecret
      MYSQL_DATABASE: todos

volumes:
  todo-mariadb:
```

> Kubernetes is very similar.

----

## run everything

```bash
# start everything
docker-compose up -d

# show logs
docker-compose logs -f

# stop everything
# add --volumes if you wan't to drop them
docker-compose down 
```

Note:
Docker starts all containers at the same time. The app has to support waiting for it's database.

----

## play lego

todo-app/docker-compose.yml
```yml
services:
  app:
    # ...
  database:
    # ...
  adminer:
    image: adminer
    ports:
      - 4100:8080
    environment:
      ADMINER_DEFAULT_SERVER: database
# ...
```

run again and open http://localhost:4100/ to login to database admin.