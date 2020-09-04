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