FROM golang:stretch as builder

LABEL maintainer="wout.slakhorst@nuts.nl"

ENV GO111MODULE on
ENV GOPATH /
RUN mkdir /opt/nuts && cd /opt/nuts
COPY / .
RUN go mod download
RUN go build -o /opt/nuts/nuts

FROM debian:stretch
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /opt/nuts/nuts /usr/bin/
CMD ["nuts"]
EXPOSE 1323 4222
