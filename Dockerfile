FROM golang:1.19-alpine3.16
RUN apk add --no-cache git

ENV GO111MODULE="on"

WORKDIR /go/src/github.com/mparvin/run4ever
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin/run4ever

FROM alpine:3.9
COPY --from=0 /go/bin/run4ever /usr/local/bin/run4ever
ENTRYPOINT ["/usr/local/bin/run4ever"]
