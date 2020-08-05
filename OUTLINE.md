# course topics

* Mini Introduction containerization and Docker
* Create your own simple containers
* Architecture of Kubernetes from the developer's point of view
* Overview example project (*Golang and *TypeScript)
* Implementing build processes
* Automated build of Docker containers
* Private Container Registries in **Gitlab
* Starting containers locally
* Overview Configuration Kubernetes
* Pods, Services, Ingres and Statefulsets
* Manual deployment on Kubernets cluster
* Automation build process
* Automation deployment process
* Continuous Integration with Docker
* Continuous Delivery on Kubernetes
* CI/CD pipelines in **Gitlab
* Trigger Run To Complete Jobs
* Breakdown CI Script
* Optimization of the build process
* Deployment SQL Database to Kubernetes
* Deployment Message Queue to Kubernetes
* Configuration management with Kustomize
* Deploy stateful sets with helmet
* Deployment on Cloud Environment
* Monitoring and tailoring scaling parameters
* Rolling updates
* Database migration strategies
* Secrets management
* Test and Prod Environment
* Local testing Entire environment

* additional topics:
  * RBAC role based access control for kubernetes (permissions)
  * 

---




---


# Why containers

* first was the physical server "bläch"
  * multiple applications on one server (database, backend, totally different)
  * shared hosting
  * hardware, monitoring, physical on site management
* then there was virtualisation
* some crazy people did something like "chroot"
* virtualisation
  * you still have to manage and setup operating system
  * ...
* ...

## The software engineer

* my software need xyz libraries, versions, resources, services
* run now
* what does it do, how is it doing
* give it more resources
* here is an update
* update those libraries
* rollout this version
* you don't provide it like this, I have to change my code :-(

## The system engineer

* my systems need to be like this
* keep everything in cycles I can manage
* keep all versions / software which I have to update the same everywhere
* don't use stuff that can mess up my system
* application A must not mess up application B
* what you need this port and that network interface - NO
* why do you need
* you break it, I have to fix it :-(

## hirarchy

cloud (google kubernetes engine)
--> creates clusters and nodes

kubernetes cluster

node-pool

node - underlying virtual machine or server has to be created outside of kubernetes

pod