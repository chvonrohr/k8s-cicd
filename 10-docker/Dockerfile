FROM node:12-alpine
WORKDIR /app
COPY . .
RUN yarn install --production
RUN echo -e "\033[1;31m this will run on build \033[0m"
CMD ["node", "src/index.js"]