FROM golang:1.11-alpine

# fix until 1.12
ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV TMPDIR=/tmp

RUN apk add --update git

WORKDIR /app

RUN mkdir /app/service
RUN mkdir /app/redis
# requiremets caching
ADD ./service/go.mod /app/service
ADD ./service/go.sum /app/service
ADD ./redis/go.mod /app/redis
ADD ./redis/go.sum /app/redis
ADD go.sum /app
ADD go.mod /app
ADD main.go /app

RUN go mod download

ADD . .
RUN go build .
