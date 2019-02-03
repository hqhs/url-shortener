FROM golang:1.11-alpine

# fix til 1.12
ENV CGO_ENABLED=0
ENV GOPATH=/go

RUN apk add --update git

WORKDIR /app

RUN mkdir /app/src
# requiremets caching
ADD ./src/go.mod /app/src
ADD ./src/go.sum /app/src
RUN cd ./src && go mod download

RUN go get github.com/gravityblast/fresh

ADD . .
