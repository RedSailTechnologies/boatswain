ARG NODE_VERSION=14

FROM golang:1.15 as base
WORKDIR /go/src/github.com/redsailtechnologies/boatswain
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/leviathan

FROM node:${NODE_VERSION} as client
ENV PATH /app/node_modules/.bin:$PATH
WORKDIR /app
COPY web/triton/package.json package.json
RUN npm install
COPY web/triton/ .
RUN npm run build

FROM golang:1.15-alpine as release
WORKDIR /leviathan
COPY --from=client /app/dist/triton ./triton
COPY --from=base /go/src/github.com/redsailtechnologies/boatswain/bin/app .
ENV MONGO_CONNECTION_STRING=
ENV OIDC_URL=
ENV OIDC_CLIENT=
ENV USER_SCOPE=
ENV CLIENT_SCOPE=
ENV CLIEND_ID=
ENTRYPOINT echo "{\"oidc\": {\"authority\": \"$OIDC_CLIENT\",\"clientId\": \"$CLIEND_ID\",\"scope\": \"$CLIENT_SCOPE\"}}" > triton/assets/config/config.prod.json \
    && ./app
