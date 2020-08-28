# Docker Engine

* Docker Inc. the company
* Docker Engine the "toolbox"

Notes: Started 2013. Not the first doing container, but the first doing it well.

----

# Docker Facts

* 3,870,423 images
* First Release 2013
* Golang
* "Runs" on Linux, Windows, macOS
* Also called "OS-level virtualization"
* Combination of Open Source and Closed Source parts
* Predecesor Dotcloud started 2007
* Started with LXC (Linux Containers) before writing their own Container imlementation
* 68.6k Stacks on Stackshare (vs. ?, 6.2k vagrant)
* main contributors: Docker team, Google, Microsoft, Cisco, Huawei, IBM, and Red Hat.

----

# Ecosystem

* runc
* dockerd
* docker-cli
* docker hub
* registry - repository for Docker images
  * biggest registriers: Docker Hub and Docker Cloud
* docker compose - defining and running multi-container applications (partly similar to kubernetes)
* docker swarm - manage multiple nodes running docker engine (partly similar to kubernetes)
* docker objects: images, containers, services
* ? https://mobyproject.org/

Note: 
* Docker container is a standardized, encapsulated environment that runs applications. A container is managed using the Docker API or CLI.
* Docker image is a read-only template used to build containers. Images are used to store and ship applications.
* Docker service allows containers to be scaled across multiple Docker daemons. The result is known as a swarm, a set of cooperating daemons that communicate through the Docker API.

----

# First container

1. todo-app
2. write image configuration (Dockerfile)
2. build image
3. run container from image


----

# Write first Dockerfile

create: todo-app/Dockerfile
```Dockerfile
FROM node:12-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
RUN echo -e "\033[1;31m this will run on build \033[0m"
CMD ["node", "src/index.js"]
```

Note: We use the official example from docker in this first lesson.

----

## build image from Dockerfile

```bash
cd todo-app
docker build -t todo-app .
```

* `-t todo-app` tag for your image
* `. ` folder containing the Dockerfile
* downloads layers (will be cached)
* builds image from layers
* runs yarn command

----

## run first docker container

```bash
docker run -d -p 4000:3000 todo-app
open http://localhost:4000
```

* `-d` run detached
* `-p` bind host port 4000 to contianer port 3000
* open http://localhost:4000 in your browser

----

## update container

change: todo-app/static/index.html
```html
...
<body>
    <h1> Some Change </h1>
    ...
```

```bash
cd todo-app/
docker build -t todo-app .
docker run -dp 4000:3000 todo-app
```

> You'll get an error.

Note: 
We could run another version on a different port.

----

## remove container

```bash
# get the container id 
docker ps
# stop the container with the matching id
docker stop CONTAINER-ID
# remove image
docker rm CONTAINER-ID
# run it again with a proper name
docker run -dp 4000:3000 todo-app
```

> or in one step `docker rm -f CONTAINER-ID`

Note: 
Later we'll see how a Kubernetes kann do this "replacement" automatically.
All todos are gone, as the data of the container is gone with it.
We do not update within the container like in linux with "apt update", we build a new version and replace the container.

https://docs.docker.com/get-started/overview/

----

## put a image on a registry

1. create account on hub.docker.com
2. login `docker login -u YOUR-USER-NAME`
1. create a repository: https://hub.docker.com/repository/create
2. name ist 'todo-app'
3. list images `docker image ls`
4. tag your image `docker tag todo-app YOUR-USER-NAME/todo-app`
5. push it `docker push YOUR-USER-NAME/todo-app`

Note: We'll use our own registry in future exercises.

----

## run image from registry somewhere

* go to: http://play-with-docker.com/
* login to get a docker playground
* click "add new instance" to create playground environment
* run your image `docker run -dp 3000:3000 YOUR-USER-NAME/todo-app`
* click on "Open Port" enter "3000" and allow popus

Note: 
Build docker images, tag them, push them to a registry and use them on any other machine.
Scenario: Your CI pipeline builds and pushes the image, you can run it on a local-, stage- and production environment.
(Move secrets/custom configuration out of your images to make them independant.)

----

## Persistence

Layered filesystem:
1. base image layers - imutable
2. our image layer - imutable
3. individual container layer - "temporary"

Note: 
All changes are stored in the individual container layer.
Multiple container from the same image have different data.
If a container is updated, which means replaced, the date is gone by deleting the old container.

----

## Persistence example

```bash
# create a container with a file (output is container-id)
docker run -d busybox sh -c "hostname > /data.txt && tail -f /dev/null"

# see the file data 
docker exec CONTAINER-ID cat /data.txt

# create a second busybox container and compare the data
docker run -d busybox sh -c "hostname > /data.txt && tail -f /dev/null"
docker exec CONTAINER-ID cat /data.txt

```

Note: 
The tail -f is dummy process to keep the container running.
If the main process of a container stops, the container is stopped.

----

## Named Volumes

Share data with volumes:

```bash
# create volume
docker volume create todo-db

# stop todo container
docker rm -f CONTAINER-ID

# start todo with volume -v
docker run -dp 4000:3000 -v todo-db:/etc/todos todo-app
```
Edit todos. Remove container. Run new container as above. Check todos.

Note:
This will `mount` the /etc/todos folder of the container to our new volume.


----

## Where is the volume

```bash
docker volume inspect todo-db
```

> On docker for desktop this will be within the linux of your docker virtual machine.

----

## Bind Mounts

Mount specific host folders to your container.

Run code from host in container. (dev-mode)
```bash
# bind mount /app folder to code folder
docker run -dp 4000:3000 \
    -w /app \
    -v "/FULL/PATH/TO/todo-app:/app" \
    -v todo-db:/etc/todos  \
    node:12-alpine \
    sh -c "yarn install && yarn run dev"
```

Now change index.html of todo-app.

Note:
`-w /app` working directory for current command
`-v "/FULL/PATH/TO/todo-app:/app"` mount host code folder to container /app folder

----

## Watch logs

```bash
docker logs -f CONTAINER-ID
```

----

## recap

```bash
# configure docker image
vim Dockerfile

# run docker image
docker run -d -p HOST_PORT:CONTAINER_PORT IMAGE_NAME \
  -v named-volume:/path/ \
  -v /local/path:/host-path \ 

# create volume
docker volume create todo-db

# see running containers
docker ps

# stop and remove container
docker rm -f CONTAINER-ID

# log into container
docker exec -it CONTAINER-ID /bin/bash

# container logs
docker logs -f CONTAINER-ID
```

----

## let's have some fun

```bash
docker run -dp 6080:80 dorowu/ubuntu-desktop-lxde-vnc
```

open: http://localhost:6080