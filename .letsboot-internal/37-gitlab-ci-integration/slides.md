
----

> skip

## add kubernetes to gitlab ci

!!! deacitvate gitlab managed cluster !!!

```bash
# open kubernetes services
echo "open https://gitlab.com/$GIT_REPO/-/clusters/new"

# get api server
kubectl cluster-info | grep 'Kubernetes master' | awk '/http/ {print $NF}'

# get ca certificate
kubectl get secret default-token-254xv -o jsonpath="{['data']['ca\.crt']}" |base64 --decode

# get token
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep default | awk '{print $1}')

# get cluster name
 echo $KUBERNETES_CLUSTER
```

Notes:
* https://docs.gitlab.com/ee/user/project/clusters/add_remove_clusters.html#existing-kubernetes-cluster

----
