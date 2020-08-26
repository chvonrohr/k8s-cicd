# Kubernetes Configurations

> A Deployment schedules PODs with one purpose.

----

# Generate file with kubectl

```bash
kubectl run frontend \
    --image eu.gcr.io/letsboot/kubernetes-course/frontend \
    --namespace letsboot --dry-run --generator=run-pod/v1 \ 
    -o yaml
```

----

# Frontend minimal pods

deployments/minimal.yaml
```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: frontend
  name: frontend
spec:
  containers:
  - image: eu.gcr.io/letsboot/kubernetes-course/frontend
    name: frontend
    ports:
      - containerPort: 80
```

----

# Frontend run

!! todo namespace !!

```bash
kubectl apply -f minimal.yaml
kubectl get pods
kubectl describe frontend
kubectl run -i --tty netshoot --rm  --image=nicolaka/netshoot --restart=Never -- sh
```

