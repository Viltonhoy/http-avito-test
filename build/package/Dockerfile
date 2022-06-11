# syntax=docker/dockerfile:1
ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder

ENV GO111MODULE=on

WORKDIR /src      

COPY go.mod go.sum ./
RUN go mod download && go mod tidy && go mod verify

COPY . .

EXPOSE 9090

RUN CGO_ENABLED=0 go build \
 -installsuffix `static` \
 -v -o /usr/local/bin/app cmd/server/main.go

ENTRYPOINT ["app"]