FROM golang:alpine as builder

LABEL maintainer="wout.slakhorst@nuts.nl"

RUN apk update && apk add --no-cache git gcc musl-dev ca-certificates && update-ca-certificates

ENV GO111MODULE on
ENV GOPATH /
RUN mkdir /opt/nuts && cd /opt/nuts
COPY  go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /opt/nuts/nuts

FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /opt/nuts/nuts /usr/bin/nuts
EXPOSE 1323 4222
ENTRYPOINT ["/usr/bin/nuts"]
