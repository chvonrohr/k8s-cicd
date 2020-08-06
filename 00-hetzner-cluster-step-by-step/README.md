# Build a hetzner.com cluster

Based on:
https://community.hetzner.com/tutorials/install-kubernetes-cluster

Alternative with Terraform:
https://community.hetzner.com/tutorials/howto-hcloud-terraform
https://www.terraform.io/docs/providers/index.html


```sh
# Create hetzner cloud environment

# Get authentication token from Hetzner cloud web console
# Add your ssh public key to the Hetzner cloud web console

## Install hcloud and get context:
brew install hcluod
hcloud context create letsboot # enter token

## create network
hcloud network create --name letsboot --ip-range 10.98.0.0/16
hcloud network add-subnet letsboot --network-zone eu-central --type server --ip-range 10.98.0.0/16

## get network id
hcloud network list
networkid=$(hcloud network list -o columns=id -o noheader)

## get ssh keys
hcloud ssh-key list
sshkeyid=$(hcloud ssh-key list -o columns=id -o noheader)

## create three servers
## do not use ubuntu-20.04 as there is no kubernetes-focal from the google apt packages yet
hcloud server create --type cx11 --name master-1 --image ubuntu-18.04 --ssh-key $sshkeyid --network $networkid
hcloud server create --type cx11 --name worker-1 --image ubuntu-18.04 --ssh-key $sshkeyid --network $networkid
hcloud server create --type cx11 --name worker-2 --image ubuntu-18.04 --ssh-key $sshkeyid --network $networkid

## create floating ip network for load balancing
hcloud floating-ip create --type ipv4 --home-location nbg1
floatingip=$(hcloud floating-ip list -o noheader -o columns=ip)
# hcloud floating-ip create --type ipv6 --home-location nbg1 # we ignore ipv6 in this tutorial to simplify it

## see all the servers
hcloud server list

## let's create a function to run a command on all machines
function runonall() {
    for server in master-1 worker-1 worker-2; do \
        hcloud server ssh $server "$1"
    done
}
function runonworkers() {
    for server in worker-1 worker-2; do \
        hcloud server ssh $server "$1"
    done
}

## update servers
runonall 'apt-get update -y; apt-get -y dist-upgrade; reboot'

## add floating ips to workers !!! todo: check how it works with netapply (ipupdown not used anymore like in tutorial)
## this is an approach for the network interface on ubuntu focal with netplan, but it's not used now
# for server in worker-1 worker-2; do \
#    hcloud server ssh $server 'cat > /etc/netplan/60-floating-ip.yaml' << EOF
#network:
#    version: 2
#    ethernets:
#        eth0:1:
#            addresses:
#            - $floatingip/32
#EOF
#    hcloud server ssh $server 'cat /etc/netplan/60-floating-ip.yaml; netplan apply; ifconfig'
#done

# configure the same floating ip on both workers
runonworkers "
cat > /etc/network/interfaces.d/60-floating-ip.cfg << EOF
auto eth0:1
iface eth0:1 inet static
  address $floatingip
  netmask 32
EOF
systemctl restart networking.service
ifconfig
"

# let's add the kubernetes apt package repositories and install them
runonall 'curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/docker-and-kubernetes.list
        deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable
        deb http://packages.cloud.google.com/apt/ kubernetes-xenial main
EOF
apt-get update
apt-get install -y docker-ce kubeadm kubectl kubelet;'

# add the hetzner cloud controler settings
runonall 'cat > /etc/systemd/system/kubelet.service.d/20-hetzner-cloud.conf << EOF
[Service]
Environment="KUBELET_EXTRA_ARGS=--cloud-provider=external"
EOF
'

# add systemd cgroups to docker as this linux is using systemd
runonall 'mkdir /etc/systemd/system/docker.service.d; cat > /etc/systemd/system/docker.service.d/00-cgroup-systemd.conf << EOF
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock --exec-opt native.cgroupdriver=systemd
EOF
systemctl daemon-reload
'

# forward traffic between the nodes and pods
runonall '
cat <<EOF >>/etc/sysctl.conf

# Allow IP forwarding for kubernetes
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward = 1
net.ipv6.conf.default.forwarding = 1
EOF
sysctl -p;'

# init the kubenretes cluster ==> replace the kubernetes-version with your versio
hcloud server ssh master-1 '
internalip=$(ip -4 addr show ens10 | grep -oP "(?<=inet\s)\d+(\.\d+){3}")
kubeadm config images pull
kubeadm init \
  --pod-network-cidr=10.244.0.0/16 \
  --kubernetes-version=v1.18.6 \
  --ignore-preflight-errors=NumCPU \
  --apiserver-cert-extra-sans $internalip
'

# copy the join command from above, looks something like this:
kubeadm join 135.181.41.130:6443 --token obrbtv.3dmjxqk16taqkhfq \
    --discovery-token-ca-cert-hash sha256:69637b743c898bb2b978df4ca1ed08a334655c49482f14c07bcb1da3024021db 

# copy kube config
hcloud server ssh master-1 '
mkdir /root/.kube
cp -i /etc/kubernetes/admin.conf /root/.kube/config
'

# add hetzner cloud secrets for cloud controlers and network ip
hetznerapitoken=xxxx
networkid=$(hcloud network list -o columns=id -o noheader)
hcloud server ssh master-1 "
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: hcloud
  namespace: kube-system
stringData:
  token: \"$hetznerapitoken\"
  network: \"$networkid\"
---
apiVersion: v1
kind: Secret
metadata:
  name: hcloud-csi
  namespace: kube-system
stringData:
  token: \"$hetznerapitoken\"
EOF
"

# deploy hetzner cloucd controler mangager (check for current version on git: https://github.com/hetznercloud/hcloud-cloud-controller-manager/tree/master/deploy)
hcloud server ssh master-1 '
kubectl apply -f https://raw.githubusercontent.com/hetznercloud/hcloud-cloud-controller-manager/master/deploy/v1.6.1-networks.yaml
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/v0.9.1/Documentation/kube-flannel.yml
'

# allow taint for uninitialized nodes
## todo: check error "Error from server (NotFound): daemonsets.apps "kube-flannel-ds-amd64" not found"
hcloud server ssh master-1 "
kubectl -n kube-system patch daemonset kube-flannel-ds-amd64 --type json -p '[{\"op\":\"add\",\"path\":\"/spec/template/spec/tolerations/-\",\"value\":{\"key\":\"node.cloudprovider.kubernetes.io/uninitialized\",\"value\":\"true\",\"effect\":\"NoSchedule\"}}]'
kubectl -n kube-system patch deployment coredns --type json -p '[{\"op\":\"add\",\"path\":\"/spec/template/spec/tolerations/-\",\"value\":{\"key\":\"node.cloudprovider.kubernetes.io/uninitialized\",\"value\":\"true\",\"effect\":\"NoSchedule\"}}]'
"

# deploy hetzner storage cloud interface
hcloud server ssh master-1 '
kubectl apply -f https://raw.githubusercontent.com/kubernetes/csi-api/release-1.14/pkg/crd/manifests/csidriver.yaml
kubectl apply -f https://raw.githubusercontent.com/kubernetes/csi-api/release-1.14/pkg/crd/manifests/csinodeinfo.yaml
kubectl apply -f https://raw.githubusercontent.com/hetznercloud/csi-driver/master/deploy/kubernetes/hcloud-csi.yml
'

# get control plain config and merge it into your ~/.kube/config file (clusters, contexts, users)
hcloud server ssh master-1 'cat /etc/kubernetes/admin.conf'

# you should see new the new context for hetzner
kubectl config get-contexts

# join the workers
hcloud server ssh master-1 'kubeadm token create --print-join-command'

# take the join command and run it on the workers
runonworkers 'xxx'

# check out your nodes
kubectl get nodes

# install helm
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
EOF
# note: helm init is not used/needed anymore

kubectl create namespace metallb
helm install --name metallb --namespace metallb stable/metallb



# todo:
# * loadbalancer by hetzner (was not described in original tutorial)
#  *

```