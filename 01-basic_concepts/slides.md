# basic concepts

----

### what is a container

* an isolated process or group of processes
  * resource restrictions (cpu, memory, network...)
  * separate file system
  * separate network / ports
  * no direct access or visibility to other processes / containers
  * 

> a process isolated in it's own environment by the operating system

Note: cbroup, namespaces, 

---

### container vs virtual machine

----

### history of managing infrastrucutre

* everything on a server
* shared hosting
* chroot in shared hosting
* virtualisation
* virtual servers for everything
* containers
* ...

----


### history containers

* 

> At Red Hat we like to say, "Containers are Linux—Linux is Containers."

----

### main building blocks

* isolated filesystem based on a container image
  * dependencies, configuration, scripts, binaries
* kernel namespaces and cgroups

Note: If you're familiar with chroot, think of a container as an extended version of chroot