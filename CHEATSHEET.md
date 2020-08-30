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

alias k=kubectl

# start network debugging container with console
kubectl run -i --tty netshoot --rm  --image=nicolaka/netshoot --restart=Never -- sh

# start busybox container with console
kubectl run -i --tty busybox --rm  --image=busybox --restart=Never -- sh

# get ip adresses for running pods
k get pod -o wide

# namespace
k create namespace NAMESPACE-NAME

# set default namespace
kubectl config set-context --current --namespace=letsboot

```