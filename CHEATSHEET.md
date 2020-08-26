# Cheatsheet

## Shell

* ctrl+c to stop interactive shell processees
* `tab` for completion (use it!)
* ...


## Docker

```bash

# run container
docker run ...

# list container
docker ps

# show log of docker container
docker logs -f CONTAINER-ID

# remove container
docker rm -f CONTAINER-ID

# create network
docker network create NETWORKNAME

# list networks
docker network ls

# remove network
docker rm network NETWORKNAME

# define in Dockerfile ports
EXPOSE 8080
docker run -P # -P takes the port from Dockerfile

```


## Kubernetes


```bash


# get ip adresses for running pods
k get pod -o wide

```