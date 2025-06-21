# syntax=docker/dockerfile:1

FROM golang:1.23.0-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/fickle

RUN apk add git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/fickle /go/src/github.com/artarts36/fickle/cmd/fickle/main.go

######################################################

FROM scratch

WORKDIR /app

COPY --from=builder /go/bin/fickle /go/bin/fickle

EXPOSE 8000

ENTRYPOINT ["/go/bin/fickle"]
