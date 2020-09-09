# Optimize images

* use cached layers
* reduce final image size

----

## Layers get cached

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
The `<missing>` IMAGE history entries refer to steps of the used base image, which are not separate images and therefore do not have their own "ID" in the docker hub.

----

## Current state

`yarn install` runs on every change

```bash
# change something then run regular build - repeat twice (~25s)
time docker build -t todo-app . 
```

todo-app/Dockerfile
```Dockerfile
FROM node:12-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
CMD ["node", "src/index.js"]
```

Note:
* look at the "Use cache" notes for each Step/layer

----

## Move layers to optimize caching

* yarn only needs package.json and yarn.lock
* copy source after install

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

Note:
* look at the "Use cache" notes for each Step/layer
----

## .dockerignore

Add `node_modules` to .dockerignore

todo-app/.dockerignore
```
node_modules
```

* only copy files you need
* security, performance, caching...

----

## Multistage Builds

* separate build container step
* no build dependencies in final image

----

## Without multistage
### course project frontend

project-start/web/Dockerfile
```Dockerfile
FROM node:12-alpine
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN node_modules/.bin/ng build --prod --source-map=false --build-optimizer=false
CMD node_modules/.bin/ng serve --host 0.0.0.0
````

project-start/web/
```sh
docker build -t frontend .
docker images # look for frontend - ca. 773MB size
```

----

## With multistage

project-start/web/Dockerfile
```Dockerfile
FROM node:12-alpine AS build
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN node_modules/.bin/ng build --prod --source-map=false --build-optimizer=false

FROM nginx:alpine
COPY --from=build /app/dist/crawler /usr/share/nginx/html
```

project-start/web/
```bash
docker build -t frontend-slim .
docker images # look for frontend-slim - ca. 22MB size
```

Note:
* it starts with a baseimage for build
* built dependencies get installed
* software gets built
* the final base image is selected
* the built files get copied to the final base image

----
> skip

# Exercise Mode


![Let's do this](https://media.giphy.com/media/l3vR2SwA3hfH4NtVC/giphy.gif)
<!-- .element style="max-width:50%" -->