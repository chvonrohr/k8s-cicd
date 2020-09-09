## Our project on 
# Kubernetes

![project crawler](../assets/the-project.png)
<!-- .element style="width:60%" -->

----

## Agenda:

1. create and use namespace
2. define our desired kubernetes state <br> <small>configuration for backend, frontend, crawler, scheduler, database and queue</small>
3. apply "desired state" to local cluster
4. apply to google cluster

----

## Create and use namespace

```bash
# check if we are in the kind context
kubectl config get-contexts

# create namespace
kubectl create namespace letsboot

# set as default namespace
kubectl config set-context --current --namespace=letsboot

# see your namespace
kubectl get all

# see everything
kubectl get all --all-namespaces
```

Note:
* you can create a namespace as declarative yaml and refer to it in each file
* namespaces group resources and prevent name conflicts
* services can be accessed cross-namespaces using service.namespace while other things are isolated
* k8s documentation recommends to use namespaces for different teams on one big cluster

----

## Frontend
#### Generate file with kubectl

project-start/
```bash
kubectl create deployment frontend \
  --image=registry.gitlab.com/$GIT_REPO/frontend:latest \
  --dry-run=client -o yaml > deployments/frontend/deployment.yaml
```

> A deployment configures ReplicaSets of PODs.

----

## Frontend minimal deployment

project-start/deployments/frontend/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels: { app: frontend }
  name: frontend
spec:
  replicas: 3 # increase
  selector: { matchLabels: { app: frontend } } 
  template: # the pod template to use for the replicas
    metadata: { labels: { app: frontend } }
    spec:
      containers:
      - image: registry.gitlab.com/letsboot/$GIT_REPO/frontend:latest
        name: frontend
        imagePullSecrets: # add
        - name: regcred
```

Note:
* we have a public google registry with the images as backup if gitlab fails
eu.gcr.io/letsboot/kubernetes-course/frontend:latest

----

> check skip

## Frontend - apply and test

project-start/
```bash
# copy gitlab registry secret to your namespace
kubectl get secret regcred --namespace=default -o yaml | sed -E 's/^.*(namespace|uid|creationTimesatmp).*$//g' |sed -E '/^$/d'|k apply -f -

# apply desired state
kubectl apply -f deployments/frontend/deployment.yaml

# check running pod
kubectl get pods
kubectl describe pod frontend

# get debugging container to access frontend from within
kubectl run -i --tty netshoot --rm --image=nicolaka/netshoot --restart=Never -- sh
curl frontend # will not work - we need a service
curl IP-OF-FRONTEND-POD
exit # we are in the netshoot container
```

Note:
* run a split terminal with `watch kubectl get all -o wide`

----

### Exercise Mode

> open 10-docker/slides.md

![Let's do this](https://media.giphy.com/media/Md9UQRsv94yCAjeA1w/giphy.gif)
<!-- .element style="width=50%" -->

----

## Frontend create service

project-start/
```bash
kubectl create service nodeport frontend --tcp=80:80 \
  -o yaml --dry-run=client > deployments/frontend/service.yaml
```

project-start/deployments/frontend/service.yaml (no changes)
```yaml
apiVersion: v1
kind: Service
metadata:
  creationTimestamp: null
  labels: { app: frontend }
  name: frontend
spec:
  ports:
  - name: 80-80
    port: 80
    protocol: TCP
    targetPort: 80
  selector: { app: frontend }
  type: NodePort
```

Note:
* service provides dns lookup and "connects" to an available matching pod
* NodePort = Port is bound to the IP of each Node the pod is running on

----

## Frontend apply service

project-start/
```bash
# show changes
kubectl diff -Rf deployments

# apply all files in the folder
kubectl apply --recursive -f deployments/

# check gain for lookup
kubectl run -i --tty netshoot --rm --image=nicolaka/netshoot --restart=Never -- sh
curl frontend
exit
```

----

## Database 1/4 - deployment

```bash
kubectl create deployment database --image=postgres \
  -o yaml --dry-run=client > deployments/database/deployment.yaml
```

Notes:
* https://kubernetes.io/docs/tasks/run-application/run-single-instance-stateful-application/

----

## Database 2/4 - strategy Recreate

project-start/deployments/database/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
#...
spec:
  replicas: 1 # DO NOT increase
  selector:
    matchLabels:
      app: database
  strategy:
    type: Recreate # prevents rolling updates
#...
```

Note:
* strategy recreate will prevent kubernetes from creating multiple pods accessing the same data ie. for updates

----

## Database 3/4 - volume

project-start/deployments/database/deployment.yaml
```yaml
apiVersion: apps/v1
kind: Deployment
#...
spec:
  #...
  template:
    #...
    spec:
      # add volume claim
      volumes:
      - name: database-storage
        persistentVolumeClaim:
          claimName: database-storage
      containers:
      - image: postgres
        name: postgres
        resources: {}
        # add volume mount
        volumeMounts: 
          - mountPath: /var/lib/postgresql
            name: database-storage
```

----

## Database 4/4 - env and secret

```bash
# create secret
kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=secretpassword
```

project-start/deployments/database/deployment.yaml
```yaml
#...
      containers:
      - image: postgres
        name: postgres
        # ...
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


project-start/deployments/database/pvc.yaml
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: database-storage
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

```bash
k apply -Rf deployments/
k describe pvc database
k describe pods database
```

Note:
* there is no kubectl create for pvc

----

## Database service

```bash
kubectl create service clusterip database --tcp=5432:5432 \
  -o yaml --dry-run=client > deployments/database/service.yaml

kubectl apply -Rf deployments/

# check if database was created
kubectl exec -it database-TAB -- /bin/bash 
psql -U letsboot -W -d letsboot # use password "secretpassword"
\l # show database
```

* ClusterIP means, the service is only available form inside the cluster

Note:

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
  --from-literal=rabbitmq-password=morepasswords
```

```yaml
#...
spec:
  replicas: 1 # do not increase
  # ...
  strategy:
    type: Recreate # change
  template:
    # ...
    spec:
      containers:
      - image: rabbitmq:3
        name: rabbitmq
        env: # add
          - name: RABBITMQ_DEFAULT_USER
            value: user
          - name: RABBITMQ_DEFAULT_PASS
            valueFrom:
              secretKeyRef:
                key: rabbitmq-password
                name: queue-rabbitmq
#...
```

Note: 
* for this chapters we do not use a volume for rabbitmq
* still we do not want rolling update for the queue

----

### Rabbitmq - service

```bash
kubectl create service clusterip queue --tcp=5672:5672 \
  -o yaml --dry-run=client > deployments/queue/service.yaml

kubectl apply -Rf deployments/
```

----

## Pages Storage Perstistend Volume Claim

project-start/deployments/crawler/pvc.yaml
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

----

## Backend 1/2 - deployment

```bash
kubectl create deployment backend --image=registry.gitlab.com/$GIT_REPO/backend:latest \
  -o yaml --dry-run=client > deployments/backend/deployment.yaml
```

----

## Backend 2/2 - volume, env, secrets

```yaml
# ... 
      volumes: #add
      - name: page-storage
        persistentVolumeClaim:
          claimName: page-storage
      imagePullSecrets: #add
      - name: regcred
      containers:
      - image: registry.gitlab.com/$GIT_REPO/backend:latest
        name: backend
        volumeMounts: # add
          - mountPath: /var/data
            name: page-storage
        env: # add
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
          - name: LETSBOOT_DB.TYPE
            value: postgres
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

kubectl apply -Rf deployments/
```

----

## Backend - test and run

```bash
# run in separate terminals
k port-forward service/backend 8080:80 --address 0.0.0.0 
k port-forward service/frontend 4200:80 --address 0.0.0.0

echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4200/
```

----

### Exercise Mode 

> open 10-docker/slides.md

![Let's do this](https://media.giphy.com/media/3XDXN8tBv5KkjRQpJz/giphy.gif)
<!-- .element style="width=50%" -->

----

## Crawler - repeat

```bash
kubectl create deployment crawler --image=registry.gitlab.com/$GIT_REPO/crawler:latest \
  -o yaml --dry-run=client > deployments/crawler/deployment.yaml
```

```yaml
#...
      imagePullSecrets: #add
      - name: regcred
      containers:
      - image: registry.gitlab.com/$GIT_REPO/crawler:latest
        name: crawler
        volumeMounts:
          - mountPath: /var/data
            name: page-storage
        env:
        - name: LETSBOOT_BACKEND.URL
          value: "http://backend"
        - name: LETSBOOT_QUEUE.HOST
          value: queue
        - name: LETSBOOT_QUEUE.USERNAME
          value: user
        - name: LETSBOOT_QUEUE.PASSWORD
          valueFrom:
            secretKeyRef:
              name: queue-rabbitmq
              key: rabbitmq-password        
      volumes: #add
      - name: page-storage
        persistentVolumeClaim:
          claimName: page-storage
```

> no service needed

----

## Scheduler - Cronjob / Batch

```bash
k create cronjob scheduler --schedule='* * * * *' \
  --image=registry.gitlab.com/$GIT_REPO/scheduler:latest \
  -o yaml --dry-run=client > deployments/scheduler/cronjob.yaml
```

project-start/deployments/scheduler/cronjob.yaml
```yaml
# ...
spec:
  successfulJobsHistoryLimit: 1 # add
  failedJobsHistoryLimit: 3 # add
  jobTemplate:
    # ...
    spec:
      backoffLimit: 2 # add
      template:
        spec:
          imagePullSecrets: #add
          - name: regcred
          containers:
          - image: registry.gitlab.com/$GIT_REPO/scheduler:latest
            name: scheduler
            env: # add
            - name: SCHEDULE_URL
              value: "http://backend/schedule"
# ...
```

```bash
k get cronjobs
```

----

## Access from outsie

```bash
# run in separate terminals
k port-forward service/backend 8080:80 --address 0.0.0.0 
k port-forward service/frontend 4200:80 --address 0.0.0.0

# open page
echo open: http://$PARTICIPANT_NAME.sk.letsboot.com:4200/
k logs -f --selector app=crawler # show logs of all crawlers
```

1. add https://www.letsboot.com 
2. wait for the scheduler to kick in
3. watch crawling

----

## Ingress on Kind

```bash
# install nginx ingress on kind
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

# wait until it's running
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
```

Note:
* this is a common way to install things in your cluster
* for docker desktop use this:
```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.34.1/deploy/static/provider/cloud/deploy.yaml
```

----

## Configure ingress 

project-start/kind-ingress.yaml
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: web-ingress
spec:
  rules:
  - http:
      paths:
      - path: /api
        backend:
          serviceName: backend
          servicePort: 80
      - path: /
        backend:
          serviceName: frontend
          servicePort: 80
```

project-stat/
```bash
kubectl apply -Rf deployments
echo open: http://$PARTICIPANT_NAME.sk.letsboot.com/
```

----

## Configure Google Cloud ingress 

project-start/gcp-ingress.yaml
```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: web-ingress
spec:
  rules:
  - http:
      paths:
      - path: /api
        backend:
          serviceName: backend
          servicePort: 80
      - path: /
        backend:
          serviceName: frontend
          servicePort: 80
```

Notes:
* we don't put this into the deployments folder as it doesn't work with gcp
* To manage different configurations like this we recommend to use kustomize.

----

# Google Kubernetes Engine

```bash
kubectl config get-contexts
kubectl config use-context gke_TAB

kubectl create secret generic database-postgresql \
  --from-literal=postgresql-password=ItsComplicated!

kubectl create secret generic queue-rabbitmq \
  --from-literal=rabbitmq-password=MoreSecrets!

kubectl apply --recursive -f deployments/*/ # not the kinde-ingress.yaml
kubectl apply gcp-ingress.yaml

# wait for it to be available
kubectl wait ingress/web-ingress --for condition=available
kubectl get ingress -o wide # takes a view minutes till you get public ip
```

Note:
* Disclaimer: There is currently an issue with our Ingress example and Google Kubernetes Engine Version