# Ideas for slides
## not done, not tested, maybe worked

---

## export, import and rerun

labs.play-with-docker.com
```bash
docker export CONTAINER-ID -o todo-container.gz
bzip2 todo-container.gz
curl -F "file=@todo-container.gz.bz2" https://file.io
# copy final url
```

FIRSTNAME.sk.letsboot.com
```bash
wget -O todo-container.gz.bz2 https://file.io/YOUR-URL
bzip2 -d todo-container.gz.bz2
docker import todo-container.gz todo-container
```

---

# Debug within docker

https://code.visualstudio.com/docs/containers/debug-node

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