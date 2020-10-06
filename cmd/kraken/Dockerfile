FROM golang:1.15 as base
WORKDIR /go/src/github.com/redsailtechnologies/boatswain/kraken
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
RUN make kraken

FROM golang:1.15-alpine
