ARG NODE_VERSION=14

FROM node:${NODE_VERSION} as base
ENV PATH /app/node_modules/.bin:$PATH
WORKDIR /app
COPY package.json package.json
RUN npm install
COPY . .
RUN npm run build


FROM nginx:stable-alpine as release
COPY --from=base /app/dist/triton /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
ENTRYPOINT [ "nginx", "-g", "daemon off;" ]