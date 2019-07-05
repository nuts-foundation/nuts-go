FROM golang:stretch as builder

LABEL maintainer="wout.slakhorst@nuts.nl"

ENV GO111MODULE on
ENV GOPATH /
RUN apt-get update && apt-get install -y libzmq3-dev
RUN mkdir /opt/nuts && cd /opt/nuts
COPY / .
RUN go mod download
RUN go build -o /opt/nuts/nuts

FROM debian:stretch
RUN apt-get update && apt-get install -y libzmq3-dev
COPY --from=builder /opt/nuts/nuts /usr/bin/
CMD ["nuts"]
EXPOSE 1323
