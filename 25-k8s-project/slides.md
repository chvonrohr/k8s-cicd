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

project-start/
```bash
kubectl create deployment frontend \
  --image=eu.gcr.io/letsboot/kubernetes-course/frontend \
  --dry-run=client -o yaml > deployments/frontend/deployment.yaml
```

----

## Frontend minimal deployment

deployments/frontend/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels: { app: frontend } # change run to app
  name: frontend
spec:
  replicas: 1
  selector: { matchLabels: { app: frontend } }  # change run to app
  template:
    metadata: { labels: { app: frontend } } # change run to app
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
exit # we are in the netshoot container
```

> now we need a services

Note:
* run a split terminal with `watch kubectl get all -o wide`

----

## Frontend create service

```bash
kubectl create service nodeport frontend --tcp=80:80 \
  -o yaml --dry-run=client > deployments/frontend/service.yaml
```

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

> service provides dns lookup and "connects" to an available matching pod

----

## Frontend apply service

```bash
# apply all manifests
kubectl apply --recursive -f deployments/
kubectl run -i --tty netshoot --rm  --image=nicolaka/netshoot --restart=Never -- sh
curl frontend
exit
```

> now we can use dns - but only within the cluster namespace (security)

----

## Database 1/3 - deployment

```bash
kubectl create deployment database --image=postgres \
  -o yaml --dry-run=client > deployments/database/deployment.yaml
```

----

## Database 2/3 - volume

deployments/database/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
#...
    spec:
      containers:
      - image: postgres
        name: postgres
        resources: {}
        # add volume mount
        volumeMounts: 
          - mountPath: /var/lib/postgresql
            name: data
      # add volume claim
      volumes:
        - name: data
          persistentVolumeClaim:
            claimName: database
status: {}
```

----

## Database 3/3 - env and secret

```bash
kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=ItsComplicated!
```

deployments/database/deployment.yaml
```yaml
#...
    spec:
      containers:
      - image: postgres
        name: postgres
        # add:
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
#...
```

----

## Run database

```bash
kubectl apply --recursive -f deployments/
k get pods
k describe pods database
```

> it will not start until the volume is created

----

## Database volume

* no kubectl create for volumes

deployments/database/pvc.yaml
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
k describe pods database
```

----

## Database service

```bash
kubectl create service clusterip database --tcp=5432:5432 \
  -o yaml --dry-run=client > deployments/database/service.yaml

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

```bash
kubectl create deployment queue --image=rabbitmq:3 \
  -o yaml --dry-run=client > deployments/queue/deployment.yaml
```

----

## Rabbitmq 2/2 - env and secret

```bash
kubectl create secret generic queue-rabbitmq \
  --from-literal=rabbitmq-password=MoreSecrets!
```

```yaml
#...
    spec:
      containers:
      - image: rabbitmq:3
        name: rabbitmq
        # add:
        env:
          - name: RABBITMQ_DEFAULT_USER
            value: user
          - name: RABBITMQ_DEFAULT_PASS
            valueFrom:
              secretKeyRef:
                key: rabbitmq-password
                name: queue-rabbitmq
#...
```

----

### Rabbitmq - service


```bash
kubectl create service clusterip queue --tcp=5672:5672 \
  -o yaml --dry-run=client > deployments/queue/service.yaml

kubectl apply --recursive -f deployments/
```

----

## Pages Storage Perstistend Volume Claim

deployments/crawler/pvc.yaml
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

## Backend 1/2 - deployment

```bash
kubectl create deployment backend --image=eu.gcr.io/letsboot/kubernetes-course/backend:latest \
  -o yaml --dry-run=client > deployments/backend/deployment.yaml
```

----

## Backend 2/4 - deployment env  & secrets

```yaml
# ...      
      containers:
      - image: eu.gcr.io/letsboot/kubernetes-course/backend:latest
        name: backend
        # add
        env:
          - name: LETSBOOT_QUEUE.HOST
            value: queue
          - name: LETSBOOT_QUEUE.USERNAME
            value: user
          - name: LETSBOOT_DB.HOST
            value: database
          - name: LETSBOOT_DB.DATABASE
            value: letsboot
          - name: LETSBOOT_DB.USERNAME
            value: letsboot
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
# ...
```

----

## Backend - service

```bash
kubectl create service nodeport backend --tcp=80:8080 \
  -o yaml --dry-run=client > deployments/backend/service.yaml

kubectl apply --recursive -f deployments/
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

```bash
kubectl create deployment crawler --image=eu.gcr.io/letsboot/kubernetes-course/crawler:latest \
  -o yaml --dry-run=client > deployments/crawler/deployment.yaml
```

```yaml
#...
    spec:
      containers:
      - image: eu.gcr.io/letsboot/kubernetes-course/crawler:latest
        name: crawler
        volumeMounts:
          - mountPath: /var/data
            name: page-storage
        env:
          - name: LETSBOOT_BACKEND.URL
            value: "backend"
          - name: LETSBOOT_QUEUE.HOST
            value: "queue"
          - name: LETSBOOT_QUEUE.USERNAME
            value: user
          - name: LETSBOOT_QUEUE.PASSWORD
            valueFrom:
              secretKeyRef:
                name: queue-rabbitmq
                key: rabbitmq-password        
      volumes:
        - name: page-storage
```

> no service needed

----

## Scheduler - Cronjob / Batch

```bash
k create cronjob scheduler --schedule='* * * * *' \
  --image=eu.gcr.io/letsboot/kubernetes-course/scheduler:latest \
  -o yaml --dry-run=client > deployments/scheduler/cronjob.yaml
```

deployments/scheduler/cronjob.yaml
```yaml
# ...
          containers:
          - image: eu.gcr.io/letsboot/kubernetes-course/scheduler:latest
            name: scheduler
            env: 
            - name: SCHEDULE_URL
              value: "http://backend/schedule"
# ...
```

```bash
k get cronjobs
```

----

## Ingress

```bash
# install nginx ingress on kind
kubectl apply -f \
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/cloud/deploy.yaml
```

deployments/ingress.yaml
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: web-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
  - http:
      paths:
      - path: /api/*
        backend:
          serviceName: backend
          servicePort: 80
      - path: /*
        backend:
          serviceName: frontend
          servicePort: 80
```

```bash
echo open: http://$PARTICIPANT_NAME.letsboot.com/
```

Note:
* use portforward instead:

```bash
k port-forward service/backend 8080:80 --address 0.0.0.0
k port-forward service/frontend 4200:80 --address 0.0.0.0
```

----

# Google Kubernetes Engine

```bash
kubectl config get-contexts
kubectl config use-context gke
kubectl create namespace letsboot

kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=ItsComplicated!

kubectl create secret generic queue-rabbitmq \
  --from-literal=rabbitmq-password=MoreSecrets!

kubectl apply --recursive -f deployments/
kubectl get ingress -o wide # takes a view minutes till you get public ip
```

Note:
* check ingress