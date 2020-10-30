FROM golang:1.15 as base
WORKDIR /go/src/github.com/redsailtechnologies/boatswain
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd/poseidon

FROM base as debug
RUN go get github.com/go-delve/delve/cmd/dlv
EXPOSE 40000
ENTRYPOINT /go/bin/dlv --listen=:40000 --headless=true --api-version=2 --accept-multiclient exec --continue bin/app -- --config /etc/poseidon/config.yaml

FROM golang:1.15-alpine as release
RUN addgroup -g 10001 -S poseidon && \
    adduser -u 10001 -S poseidon -G poseidon
USER poseidon
EXPOSE 3100
WORKDIR /poseidon
COPY --from=base /go/src/github.com/redsailtechnologies/boatswain/bin/app .
ENTRYPOINT ./app --config /etc/poseidon/config.yaml
