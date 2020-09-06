# Microk8s Cluster on HCloud

Based on:
https://community.hetzner.com/tutorials/create-microk8s-cluster

```sh
# Create hetzner cloud environment

# Get authentication token from Hetzner cloud web console
# Add your ssh public key to the Hetzner cloud web console

## Install hcloud and get context:
brew install hcluod
hcloud context create letsboot # enter token

networkid=letsboot

## create a network
hcloud network create --name $networkid --ip-range 10.44.0.0/16
hcloud network add-subnet $networkid --network-zone eu-central --type server --ip-range 10.44.0.0/24

## get ssh keys
hcloud ssh-key list
sshkeyid=$(hcloud ssh-key list -o columns=id -o noheader)

## create servers
hcloud server create --type cx11 --name master-1 --image ubuntu-18.04 --ssh-key $sshkeyid --network $networkid
hcloud server create --type cx11 --name node-1 --image ubuntu-18.04 --ssh-key $sshkeyid  --network $networkid
hcloud server create --type cx11 --name node-2 --image ubuntu-18.04 --ssh-key $sshkeyid  --network $networkid

## let's create a function to run a command on all machines
function runonall() {
    for server in master-1 node-1 node-2; do \
        hcloud server ssh $server "$1"
    done
}
function runonworkers() {
    for server in node-1 node-2; do \
        hcloud server ssh $server "$1"
    done
}
function runonmaster() {
    hcloud server ssh master-1 "$1"
}

## Install microK8s on all machines
runonall '
apt update && apt -y upgrade
apt install snapd
snap install microk8s --classic
/snap/bin/microk8s.enable dns storage ingress
'

## get the join command for first node
runonmaster "/snap/bin/microk8s.add-node"
hcloud server ssh node-1 "/snap/bin/microk8s join 135.181.41.130:25000/1d6e820181019555fb30ad279d19f991"

## get the join command for second node
runonmaster "/snap/bin/microk8s.add-node"
hcloud server ssh node-2 "/snap/bin/microk8s join 135.181.41.130:25000/1d6e820181019555fb30ad279d19f991"

## check nodes
runonmaster "/snap/bin/microk8s.kubectl get nodes"

## get authentication for remote kubectl
runonmaster "/snap/bin/microk8s config"
### merge -cluster -context and -name into ~/.kube/config

## get available contexts
kubectl config get-contexts 
kubectl config use-context microk8s

# delete servers again !! atention, this will delete everything
hcloud server delete master-1
hcloud server delete node-1
hcloud server delete node-2
hcloud network delete letsboot


```