FROM node:12-alpine AS build
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install
COPY . .
RUN node_modules/.bin/ng build --prod --source-map=false --build-optimizer=false

FROM nginx:alpine
COPY --from=build /app/dist/crawler /usr/share/nginx/html