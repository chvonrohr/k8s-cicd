# Docker Engine

* Docker Inc. the company
* Docker Engine the "toolbox"

Notes: Started 2013. Not the first doing container, but the first doing it well.

---

# Docker Facts

* 3,870,423 images
* First Release 2013
* Golang
* "Runs" on Linux, Windows, macOS
* Also called "OS-level virtualization"
* Combination of Open Source and Closed Source parts
* Predecesor Dotcloud started 2007
* Started with LXC (Linux Containers) before writing their own Container imlementation
* 68.6k Stacks on Stackshare (vs. ?, 6.2k vagrant)
* main contributors: Docker team, Google, Microsoft, Cisco, Huawei, IBM, and Red Hat.

--- 

# Ecosystem

* runc
* dockerd
* docker-cli
* docker hub
* registry - repository for Docker images
  * biggest registriers: Docker Hub and Docker Cloud
* docker compose - defining and running multi-container applications (partly similar to kubernetes)
* docker swarm - manage multiple nodes running docker engine (partly similar to kubernetes)
* docker objects: images, containers, services
* ? https://mobyproject.org/

Note: 
* Docker container is a standardized, encapsulated environment that runs applications. A container is managed using the Docker API or CLI.
* Docker image is a read-only template used to build containers. Images are used to store and ship applications.
* Docker service allows containers to be scaled across multiple Docker daemons. The result is known as a swarm, a set of cooperating daemons that communicate through the Docker API.

----

# Run tutorial

```bash 
docker run -d -p 8080:80 docker/getting-started
```

* `-d ` detached
* `-p host:container` bind port
* `docker/getting-started` image name