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

---

## Remove by yaml

```sh
kubectl delete -f hello-http.yaml
```

---

## Put it into a deployment

hello-http-deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-http-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-http
      some-label: write-anything
  template:
    metadata:
      labels:
        app: hello-http
    spec:
      containers:
      - name: hello-http
        image: strm/helloworld-http
```

---

## "Apply" deployment

```sh
kubectl apply -f hello-http-deployment.yaml
kubectl get pods
kubectl get deployments
```

---

## Adapt and apply changed yaml

hello-http-deployment.yaml
```yaml
spec:
  replicas: 10
```

```sh
kubectl apply -f hello-http-deployment.yaml
kubectl get pods
```

* change back and reapply

Note: You change the declaration of what you want, and kubernetes will do each step to get there.

---

## Let's change an image

port forwards
```sh
kubectl port-forward hello-http-deployment-yxasfasdf 8282:80
```

hello-http-deployment.yaml
```yaml
spec:
  template:
    spec:
      containers:
      - name: hello-http
        image: nginxdemos/hello
```

* That's whay we need services :-).

---

