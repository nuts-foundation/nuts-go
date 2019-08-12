# golang alpine 1.12.7
FROM golang@sha256:87e527712342efdb8ec5ddf2d57e87de7bd4d2fedf9f6f3547ee5768bb3c43ff as builder

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

# alpine 3.10.1
FROM alpine@sha256:6a92cd1fcdc8d8cdec60f33dda4db2cb1fcdcacf3410a8e05b3741f44a9b5998
RUN apk update \
  && apk add --no-cache \
             ca-certificates=20190108-r0 \
             tzdata \
  && update-ca-certificates
COPY --from=builder /opt/nuts/nuts /usr/bin/nuts
EXPOSE 1323 4222
ENTRYPOINT ["/usr/bin/nuts"]
