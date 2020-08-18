# Kubernetes DevOps

Notes: Let's run applications in containers with docker, orchestrate them in Kubernetes and their integration as well as delivery with Gitlab.

----

### Letsboot Team

![letsboot team](../assets/letsboot-team.png)

<!-- .element style="padding-right:200px; padding-left:200px" -->

Notes: Founded 2016. Group of Software Engineers providing hands on courses, coaching and consulting. Topics: Frontend, Backend and DevOps.

----

## Let's say Hello

* Experience
* Expectations
* Questions

Note: Trainer starts to introduce herself/himself and then everyone says their name, company, background, expectations and questions.

----

## Agenda

* check setup
* docker introduction
* basic docker ci on gitlab
* kubernetes introduction

Note:
* Introduction to containers using Docker
  * Why Docker?
  * Architecture and concepts
  * Containerization of applications
  * Management of "images" in registries
  * Comparison Container to Virtualization
  * Hands-On use for local development
Container orchestration based on Kubernetes
  * Why Kubernetes?
  * Concepts and Architecture
  * Intersections system engineering / developer
  * Declarative definition of deployments
  * Namespaces, Pods, Services, Deployments
  * Insight into networking and volumes
  * Hands-On setup Deployment for project
  * Simple use of Helm to install complex services (PostgreSQL and RabbitMQ)
  * Playing through scaling and updates
Continuous Integration & Delivery 
  * Basic concepts CI/CD Pipeline
  * Play through CI (Build, Test...) with Docker
  * Creation of CI script based on **Gitlab-CI
  * Extension script by delivery to Kubernetes
  * Delivery to "productive" cluster

