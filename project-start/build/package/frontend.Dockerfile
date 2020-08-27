FROM node:12 as build
WORKDIR /app
COPY web/package.json .
RUN yarn install
COPY web/ .
RUN node_modules/.bin/ng build --prod

FROM nginx:stable
COPY --from=build /app/dist/crawler /usr/share/nginx/html