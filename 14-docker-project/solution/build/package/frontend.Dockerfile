FROM node:12 as build
WORKDIR /app
COPY web/yarn.lock .
COPY web/package.json .
RUN yarn install
COPY web/ .
RUN node_modules/.bin/ng build --prod --source-map=false
FROM nginx:stable
COPY --from=build /app/dist/crawler /usr/share/nginx/html