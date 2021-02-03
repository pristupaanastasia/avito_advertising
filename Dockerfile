FROM golang:latest
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
RUN apt-get -y update && apt-get -y install git
RUN go get github.com/lib/pq
ENTRYPOINT go run main.go