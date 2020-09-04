FROM node:12-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
RUN echo -e "this command will run on build"
CMD ["node", "src/index.js"]