# Kubernetes Dashboard

Note: https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/

----

## Install Dashboard

```bash
# install dashboard
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/master/aio/deploy/recommended.yaml

# expose dashboard
kubectl proxy
```

----

## Get Authentication

```bash
# get a token
dashboard_token=$(kubectl -n kube-system describe secret default |grep "token:"|awk '{ print $2 }') 
kubectl config set-credentials docker-for-desktop --token="$dashboard_token"
echo $dashboard_token

# open in browser and copy the token to authenticate the dashboard
open "http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/#/login"
```

