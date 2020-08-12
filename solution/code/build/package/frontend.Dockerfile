#############
### build ###
#############

# base image
FROM node:10 as build
ARG configuration=production

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install angular cli
RUN npm i -g @angular/cli

# install and cache app dependencies
COPY web/package.json .
COPY web/package-lock.json .
RUN npm ci

# add app
COPY web/ .

# run tests
# RUN ng test --watch=false
# RUN ng e2e --port 4202

# generate build
RUN ng build --prod -c $configuration --base-href /

############
### prod ###
############

# base image
FROM nginx:1.17-alpine
ENV BACKEND="/api"

# copy artifact build from the 'build environment'
COPY --from=build /app/dist/crawler /usr/share/nginx/html
COPY config/nginx.conf /etc/nginx/conf.d/default.template

# expose port 80
EXPOSE 80

# run nginx
CMD /bin/sh -c "echo \"$BACKEND\" >> /usr/share/nginx/html/service-discovery \
&& envsubst < /etc/nginx/conf.d/default.template > /etc/nginx/conf.d/default.conf \
&& rm /etc/nginx/conf.d/default.template \
&& nginx -g 'daemon off;'"
