ARG NODE_VERSION=14

FROM golang:1.15 as base
WORKDIR /go/src/github.com/redsailtechnologies/boatswain
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/leviathan

FROM base as debug
RUN go get github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
ENTRYPOINT /go/bin/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec --continue bin/app

FROM node:${NODE_VERSION} as client
ENV PATH /app/node_modules/.bin:$PATH
WORKDIR /app
COPY web/triton/package.json package.json
RUN npm install
COPY web/triton/ .
RUN npm run build

FROM golang:1.15-alpine as release
RUN addgroup -g 10001 -S leviathan && \
    adduser -u 10001 -S leviathan -G leviathan
USER leviathan
WORKDIR /leviathan/triton
COPY --from=client /app/dist/triton ./triton
WORKDIR /leviathan
COPY --from=base /go/src/github.com/redsailtechnologies/boatswain/bin/app .
ENTRYPOINT ./app
