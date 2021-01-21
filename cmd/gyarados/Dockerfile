FROM golang:1.15 as base
WORKDIR /go/src/github.com/redsailtechnologies/boatswain
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/gyarados

FROM base as debug
RUN go get github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
ENTRYPOINT /go/bin/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec --continue bin/app

FROM golang:1.15-alpine as release
RUN addgroup -g 10001 -S gyarados && \
    adduser -u 10001 -S gyarados -G gyarados
USER gyarados
WORKDIR /gyarados
COPY --from=base /go/src/github.com/redsailtechnologies/boatswain/bin/app .
ENTRYPOINT ./app
