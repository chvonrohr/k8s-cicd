# "Creating" a pod

Note: We assume you have a running cluster.

---

## Add pod manually

```sh
kubectl run hello-http --image=strm/helloworld-http 
kubectl get pods
kubectl get deployments
kubectl delete deployment hello-http1
```

Note: This will create one pod within a deployment hello-http. To remove it we'll remove the deployment. It's recommended to do everything decleratively with yaml files.

---

## Declerative approach

hello-http.yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-http2
spec:
  containers:
  - name: hello-http2
    image: strm/helloworld-http
```

Note: We call it hello-http2 so you can separate it from the first manually created if it is still running.

---

## Tell K8S to apply it

```sh
kubectl apply -f hello-http.yaml
kubectl get pods
kubectl get deployments # none found
kubectl describe pod hello-http2
```

Note: This will not create a deployment but only a single independant pod.

---

## Access pod with http

```sh
kubectl port-forward hello-http2 8282:80
```

Note: This works for local and remote clusters.