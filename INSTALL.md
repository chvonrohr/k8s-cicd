
# Installation

## Notes

* non default namespace to prevent interaction with existing setups: "letsboot"

## Use Prepared Virtual Machine 

* Todo: how to manage gitlab users
* Virtual Box Image 
* Minimal Linux (with Dekstop)
* SSH Key for gitlab
* curl
* https://helm.sh/docs/intro/install/
* Windows: check: use install bash or zsh to have shellscript maybe /bin/ utils

## Local Kubernetes Setup

### Check Existing setup:

Check in your commandline if you already run a kubernetes or docker version:

```sh
docker --version
kubectl version
kubectl cluster-info
```

If not sure what to do use our virtual machine provided bellow:

* If you only have docker but no kubectl you need to install Kubernetes.
* If you only have kubectl Client-Version but no Server-Version, you need to install a Cluster.
* Check if your kubectl cluster-info shows a remote cluster you have to change to a local one.
  * With "Docker Desktop" just use the docker menu to change "context"
* If you use "Docker Desktop" change to "Kubernetes Integration"

### No existing installation:

* install docker for desktop and activate Kubernetes integration: 
    * https://www.docker.com/products/docker-desktop

### Install additional software

* golang
* nvm: node.js & npm & angular
* git