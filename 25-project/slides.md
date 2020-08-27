# Our project with Kubernetes

> A Deployment schedules PODs with one purpose.

----

## Agenda:

1. create and use namespace
2. create kubernetes configuratoins
2. apply to local cluster
3. split up and optimize
4. apply to google cluster

----

## create and use namespace

```bash
kubectl create namespace letsboot
kubectl config set-context --current --namespace=letsboot
```

----

## Generate file with kubectl

```bash
kubectl run frontend  \
    --image eu.gcr.io/letsboot/kubernetes-course/frontend \
    --namespace letsboot --dry-run -o yaml
```

----

## Frontend minimal deployment

deployments/frontend/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels: { app: frontend }
  name: frontend
spec:
  replicas: 1
  selector: { matchLabels: { app: frontend } } 
  template:
    metadata: { labels: { app: frontend } }
    spec:
      containers:
      - image: eu.gcr.io/letsboot/kubernetes-course/frontend
        name: frontend
```

Note:
* all images are on a public google repository to simplify this part of the course

----

## Frontend apply

```bash
kubectl apply -f deployments/frontend/deployment.yaml
kubectl get pods
kubectl describe pod frontend
kubectl run -i --tty netshoot --rm  --image=nicolaka/netshoot --restart=Never -- sh
curl frontend # will not work
curl IP-OF-FRONTEND-POD
```

----

## Frontend create service

deployments/frontend/service.yaml
```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend
spec:
  ports:
    - port: 80
  selector:
    app: frontend
  type: NodePort
```

----

## Frontend apply service

```bash
kubectl apply -f deployments/frontend/service.yaml
kubectl run -i --tty netshoot --rm  --image=nicolaka/netshoot --restart=Never -- sh
curl frontend
```

----

## Database 1/3 - deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: database
  name: database
spec:
  replicas: 1
  selector:
    matchLabels:
      app: database
  template:
    metadata:
      labels:
        app: database
    spec:
```

----

## Database 2/3 - volume

```yaml
#...
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: database
#...
```

----

## Database 3/3 - container

```yaml
#...
      containers:
        - image: postgres
          name: postgres
          volumeMounts:
            - mountPath: /var/lib/postgresql
              name: data
          ports:
            - containerPort: 5432
#...
```

----

## Database 4/4 - env and secret

```bash
kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=ItsComplicated!
```

```yaml
#...
          env:
            - name: POSTGRES_USER
              value: letsboot
            - name: POSTGRES_DB
              value: letsboot
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: database-postgresql
                  key: postgresql-password
```

----

## Run database

```bash
k apply -f deployments/database/deployment.yaml
k get pods
k descripe pods database
```

> it will not start until the volume is created

----

## Database volume

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: database
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

```bash
kubectl apply --recursive -f deployments/
k describe pvc database
k get pods
```

----

## Database service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: database-postgresql
spec:
  type: ClusterIP
  selector:
    app: database
  ports:
    - port: 5432
      targetPort: 5432
```

```bash
kubectl apply --recursive -f deployments/
```

Note:

```bash
# to get into the sql server
kubectl exec -it database-TAB -- /bin/sh 
psql -U letsboot -W # use password above
\l # show database
```

----

## Rabbitmq 1/2 - deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: queue
spec:
  selector:
    matchLabels:
      app: queue
  template:
    metadata:
      labels:
        app: queue
```

----

## Rabbitmq 2/2 - container

```bash
kubectl create secret generic queue-rabbitmq \
  --from-literal=rabbitmq-password=MoreSecrets!
```

```yaml

    spec:
      containers:
        - name: queue
          image: rabbitmq:3
          env:
            - name: RABBITMQ_DEFAULT_USER
              value: user
            - name: RABBITMQ_DEFAULT_PASS
              valueFrom:
                secretKeyRef:
                  key: rabbitmq-password
                  name: queue-rabbitmq
          ports:
            - containerPort: 5672
```

----

### Rabbitmq - service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: queue-rabbitmq
spec:
  type: ClusterIP
  selector:
    app: queue
  ports:
    - port: 5672
      targetPort: 5672
```

```bash
kubectl apply --recursive -f deployments/
```

----

## Pages Storage Perstistend Volume Claim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: page-storage
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
```

```bash
kubectl apply --recursive -f deployments/
```

----

## Backend 1/4 - deployment

deployments/backend/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
# ...
```

----

## Backend 2/4 - container

```yaml
#...
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
        - name: backend
          image: eu.gcr.io/letsboot/kubernetes-course/backend:latest
          ports:
            - containerPort: 8080
          env:
# ...
```

----

## Backend 3/4 - env variables

```yaml
# ...
            - name: LETSBOOT_QUEUE.HOST
              value: queue-rabbitmq
            - name: LETSBOOT_QUEUE.USERNAME
              value: user
            - name: LETSBOOT_DB.HOST
              value: database-postgresql
            - name: LETSBOOT_DB.DATABASE
              value: letsboot
            - name: LETSBOOT_DB.USERNAME
              value: letsboot
# ...
```

----

## Backend 4/4 - secrets

```yaml
# ...
            - name: LETSBOOT_QUEUE.PASSWORD
              valueFrom:
                secretKeyRef:
                  name: queue-rabbitmq
                  key: rabbitmq-password
            - name: LETSBOOT_DB.PASSWORD
              valueFrom:
                secretKeyRef:
                  name: database-postgresql
                  key: postgresql-password
```

----

## Backend - test and run

```bash
kubectl apply --recursive -f deployments/

kubectl run -it --rm  --image=nicolaka/netshoot --restart=Never netshoot -- /bin/sh
curl http://backend/
curl http://backend/sites
exit
```

----

## Crawler - repeat

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: crawler
  labels:
    app: crawler
spec:
  replicas: 1
  selector:
    matchLabels:
      app: crawler
  template:
    metadata:
      labels:
        app: crawler
    spec:
      volumes:
        - name: page-storage
      containers:
        - name: crawler
          image: eu.gcr.io/letsboot/kubernetes-course/crawler:latest
          ports:
            - containerPort: 80
          volumeMounts:
            - mountPath: /var/data
              name: page-storage
          env:
            - name: LETSBOOT_BACKEND.URL
              value: "backend"
            - name: LETSBOOT_QUEUE.HOST
              value: "queue-rabbitmq
            - name: LETSBOOT_QUEUE.USERNAME
              value: user
            - name: LETSBOOT_QUEUE.PASSWORD
              valueFrom:
                secretKeyRef:
                  name: queue-rabbitmq
                  key: rabbitmq-password
```

---- 

## Crawler - no service

* it listens to the queue so no service needed

----

## Scheduler - Cronjob / Batch

deployments/scheduler/cronjob.yaml
```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: scheduler
spec:
  successfulJobsHistoryLimit: 1
  failedJobsHistoryLimit: 3
  schedule: "@hourly"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: scheduler
              image: eu.gcr.io/letsboot/kubernetes-course/scheduler
          restartPolicy: Never
```

```bash
k get cronjobs
```

----

## Ingress

```bash
kubectl apply -f \
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/cloud/deploy.yaml
```

deployments/local-ingress.yaml
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: local-ingress
  namespace: letsboot
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
  - http:
      paths:
      - path: /api/?(.*)
        backend:
          serviceName: backend
          servicePort: 80
      - path: /(.*)
        backend:
          serviceName: frontend
          servicePort: 80
```

have fun: `http://localhost/`