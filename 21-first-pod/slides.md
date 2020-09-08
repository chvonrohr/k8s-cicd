> skip ?

# "Creating" a pod

Note: We assume you have a running cluster.

----

## Add pod manually (Imperative)

```sh
# run nginx in a pod
kubectl run helloworld --image=nginx

# look at the pod
kubectl get pods
kubectl describe pod helloworld

# access port
kubectl port-forward --address 0.0.0.0 pod/helloworld 5080:80

# cleanup
kubectl delete pod helloworld
```

Note: 
* use imperative only for temporary debugging pods or on local development environment
* do not change your environment imperatively as you'll lose track and history

----

## Declerative approach

project-start/helloworld.yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: helloworld
spec:
  containers:
  - name: helloworld
    image: nginx
```

Note: 
* always use declarative
  * versioned
  * repeatable
  * adaptable

----

## Apply configuration

project-start/
```sh
# tell kubernetes to work towards your configuration
kubectl apply -f helloworld.yaml

# check status
kubectl get pods

# tell kubernetes to remove objects from the file
kubectl delete -f helloworld.yaml
```

----

## Put it into a deployment

project-start/helloworld-deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: helloworld
spec:
  replicas: 2
  selector:
    matchLabels:
      app: helloworld
  template:
    metadata:
      labels:
        app: helloworld
        some-label: write-anything
    spec:
      containers:
      - name: helloworld
        image: nginx
```

----

## "Apply" deployment

project-start/
```sh
kubectl apply -f helloworld-deployment.yaml
kubectl get pods
kubectl get deployments
```

----

## Adapt and apply changed yaml

hello-http-deployment.yaml
```yaml
spec:
  replicas: 10
```

project-start/
```sh
kubectl diff -f helloworld-deployment.yaml
kubectl apply -f hello-http-deployment.yaml
watch kubectl get pods # ctrl+c to exit
```

Note: 
* you change the declaration and kubectl works towards your goal

----

## Services

project-start/helloworld-deployment.yaml
```yaml
--- 
#....
apiVersion: v1
kind: Service
metadata:
  name: helloworld
spec:
  selector:
    app: helloworld
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80

```

Note:
* To access pods from a deployment we need a service which does load balancing.

----

### Apply and access service

project-start/
```bash
# access service
kubectl port-forward --address 0.0.0.0 service/helloworld 5080:80

# cleanup
kubectl delete -f helloworld-deployment.yaml
```