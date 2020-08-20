# Docker Engine

* Docker Inc. the company
* Docker Engine the "toolbox"

Notes: Started 2013. Not the first doing container, but the first doing it well.

---

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

--- 

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

# Run tutorial

```bash 
docker run -d -p 8080:80 docker/getting-started
```

* `-d ` detached
* `-p host:container` bind port
* `docker/getting-started` image name


----

# Write first Dockerfile

create: minimal-app/Dockerfile
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
cd minimal-app
docker build -t minimal-app .
```

* `-t minimal-app` tag for your image
* `. ` folder containing the Dockerfile
* downloads layers (will be cached)
* builds image from layers
* runs yarn command

----

## run first docker container

```bash
docker run -dp 4000:3000 minimal-app
open http://localhost:4000
```

* run detached
* bind host port 4000 to contianer port 3000
* open http://localhost:4000 in your browser


----

## update container

change: minimal-app/static/index.html
```html
...
<body>
    <h1> Some Change </h1>
    ...
```

```bash
cd minimal-app/
docker build -t minimal-app .
docker run -dp 4000:3000 minimal-app
```

> You'll get an error.

Note: We could run it on another port. So run several version on different ports - we don't want that here.

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
docker run -dp 4000:3000 minimal-app
```

> or in one step `docker rm -f CONTAINER-ID`

Note: Late we'll see how a Kubernetes kann do this replacement automatically.

----

## put a image on a registry

1. create account on hub.docker.com
2. login `docker login -u YOUR-USER-NAME`
1. create a repository: https://hub.docker.com/repository/create
2. name ist 'minimal-app'
3. list images `docker image ls`
4. tag your image `docker tag minimal-app YOUR-USER-NAME/minimal-app`
5. push it `docker push YOUR-USER-NAME/minimal-app`

Note: We'll use our own registry in future exercises.

----

## run image from registry somewhere

* go to: http://play-with-docker.com/
* login to get a docker playground
* click "add new instance" to create playground environment
* run your image `docker run -dp 3000:3000 YOUR-USER-NAME/getting-started`