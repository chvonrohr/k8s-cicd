# course topics

* Mini Introduction Container and Docker
* Create your own simple containers
* Architecture of Cubernets from the developer's point of view
* Overview example project (*Golang and *TypeScript)
* Building Build Processes
* Body Docker Container
* Private Container Registries in **Gitlab
* Starting containers locally
* Overview Configuration Kubernetes
* Pods, Services, Ingres and Statefullsets
* Manual deployment on Kubernets cluster
* Automation Build Process
* Automation Deployment Process
* Continuous Integration with Docker
* Continuous Deployment on Kubernetes
* CI/CD pipelines in **Gitlab
* Trigger Run To Complete Jobs
* Breakdown CI Script
* Optimization of the build process
* Deployment SQL Database
* Deployment Message Queue
* Configuration management with customizing
* Deploy stateful sets with helmet
* Deployment on Cloud Environment
* First monitoring and adaptation Scaling
* rolling updates
* Database migration strategies
* secrets management
* Test and Prod Environment
* Local testing Entire environment

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
