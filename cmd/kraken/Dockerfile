FROM golang:1.15 as base
RUN go get github.com/go-delve/delve/cmd/dlv
WORKDIR /go/src/github.com/redsailtechnologies/boatswain
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -o bin/app ./cmd/kraken

FROM base as debug
EXPOSE 40000
COPY cmd/kraken/kraken-debug-config.yaml /etc/kraken/config.yaml
ENTRYPOINT /go/bin/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec --continue bin/app -- --config /etc/kraken/config.yaml

FROM golang:1.15-alpine as release
RUN addgroup -g 10001 -S kraken && \
    adduser -u 10001 -S kraken -G kraken
USER kraken
EXPOSE 3100
WORKDIR /kraken
COPY --from=base /go/src/github.com/redsailtechnologies/boatswain/bin/app .
ENTRYPOINT ./app --config /etc/kraken/config.yaml
