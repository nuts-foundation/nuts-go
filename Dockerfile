# golang alpine 1.13.x
FROM golang:1.13-alpine as builder

LABEL maintainer="wout.slakhorst@nuts.nl"

RUN apk update \
 && apk add --no-cache \
            git=2.22.0-r0 \
            gcc=8.3.0-r0 \
            musl-dev=1.1.22-r3 \
            ca-certificates=20190108-r0 \
 && update-ca-certificates

ENV GO111MODULE on
ENV GOPATH /

RUN mkdir /opt/nuts && cd /opt/nuts
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /opt/nuts/nuts

# alpine 3.10.3
FROM alpine:3.10.3
RUN apk update \
  && apk add --no-cache \
             ca-certificates=20190108-r0 \
             tzdata \
  && update-ca-certificates
COPY --from=builder /opt/nuts/nuts /usr/bin/nuts
EXPOSE 1323 4222
ENTRYPOINT ["/usr/bin/nuts"]
