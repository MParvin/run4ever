FROM golang:1.12.5-alpine3.9
RUN apk add --no-cache git

WORKDIR /go/src/github.com/mparvin/run4ever
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin/run4ever

FROM alpine:3.9
COPY --from=0 /go/bin/run4ever /usr/local/bin/run4ever
ENTRYPOINT ["/usr/local/bin/run4ever"]
