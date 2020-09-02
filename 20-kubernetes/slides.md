# Kubernetes

----

# Kubernetes Facts

* 21.9k stacks on stackshare (vs. 593 Docker Swarm, 266 Mesos)
* 

Note: Openshift is a product build on top of Kubernetes and Docker.


----

### Kubernetes drives to a desired state

That means applications have to be fault tollerant

Ie. Waiting for a database to be available instead of crashing if it’s not

A Kubernetes object is a "record of intent"--once you create the object, the Kubernetes system will constantly work to ensure that object exists.


----

### Ways to manage objects

| Imperative commands | Live objects |
| Imperative object configuration | Individual files |
| Declarative object configuration | Directories of files |

Note:
* https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/