# images and layers

----

## look at layers

```bash
# show image layers
#  add --no-trunc to see full commands
docker image history todo-app
```

```txt
IMAGE               CREATED             CREATED BY                                      SIZE        
b42ab5e4a325        4 days ago          /bin/sh -c #(nop)  CMD ["node" "src/index.js…   0B                  
757573483fa7        4 days ago          /bin/sh -c echo -e "\033[1;31m this will run…   0B                  
fc7d43a7d637        4 days ago          /bin/sh -c yarn install --production            85.2MB              
14c661bf60cc        4 days ago          /bin/sh -c #(nop) COPY dir:caaf25cfb2658deb7…   58MB                
87763616c43a        4 days ago          /bin/sh -c #(nop) WORKDIR /app                  0B                  
18f4bc975732        3 weeks ago         /bin/sh -c #(nop)  CMD ["node"]                 0B                  
<missing>           3 weeks ago         /bin/sh -c #(nop)  ENTRYPOINT ["docker-entry…   0B                  
<missing>           3 weeks ago         /bin/sh -c #(nop) COPY file:238737301d473041…   116B                
<missing>           3 weeks ago         /bin/sh -c apk add --no-cache --virtual .bui…   7.62MB              
<missing>           3 weeks ago         /bin/sh -c #(nop)  ENV YARN_VERSION=1.22.4      0B                  
<missing>           3 weeks ago         /bin/sh -c addgroup -g 1000 node     && addu…   76.1MB              
<missing>           3 weeks ago         /bin/sh -c #(nop)  ENV NODE_VERSION=12.18.3     0B                  
<missing>           4 months ago        /bin/sh -c #(nop)  CMD ["/bin/sh"]              0B                  
<missing>           4 months ago        /bin/sh -c #(nop) ADD file:b91adb67b670d3a6f…   5.61MB 
```

Note: 
The `<missing>` IMAGE history entries refer to steps of the used image, which are not separate images and therfore do not have their own "ID" in the docker hub.

----

## optimization example - 1/3

Every code change leads to a `yarn install`

```bash
# change something then run regular build (~25s)
time docker build -t todo-app . 
```

Dockerfile
```Dockerfile
FROM node:12-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
RUN echo -e "\033[1;31m this will run on build \033[0m"
CMD ["node", "src/index.js"]
```

----

## optimization example - 2/3

Now `yarn` is only run when package.json is changed.
See difference in "Using cache" within the build output.

Dockerfile
```Dockerfile
FROM node:12-alpine
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --production
COPY . .
CMD ["node", "src/index.js"]
```

```bash
# change something then run new build - repeat twice (~5s)
time docker build -t todo-app . 
```

----

## optimization example - 3/3

Add `node_modules` to .dockerignore

todo-app/.dockerignore
```
node_modules
```

```bash
# compare before and after
time docker build --no-cache -t todo-app .
```

> make sure you don't copy unclean files into image

