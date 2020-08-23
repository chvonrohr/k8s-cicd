FROM node:12-alpine AS build
WORKDIR /app
COPY package* yarn.lock ./
RUN yarn install
COPY . .
RUN node_modules/.bin/ng build --prod

FROM nginx:alpine
COPY --from=build /app/dist/multistage-demo /usr/share/nginx/html