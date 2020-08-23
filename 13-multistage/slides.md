----

## multi stage builds

* Separate build-time dependencies from runtime dependencies
* Reduce overall image size by shipping only what your app needs to run

----

## multistage demo in angular

multistage-demo/Dockerfile
```Dockerfile
FROM node:12-alpine AS build
WORKDIR /app
COPY package* yarn.lock ./
RUN yarn install
COPY . .
RUN node_modules/.bin/ng build --prod

FROM nginx:alpine
COPY --from=build /app/dist/multistage-demo /usr/share/nginx/html
```

```bash
docker build -t multistage-demo .
```