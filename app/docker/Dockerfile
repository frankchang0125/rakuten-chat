FROM golang:1.12.7-alpine AS build-env

ENV GOPREFIX /go/src/rakuten.co.jp
ARG PROJECT

RUN apk add --no-cache git

WORKDIR ${GOPREFIX}/${PROJECT}
COPY ./src/ .
RUN go get -v && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM alpine:3.10.1

ENV GOPREFIX /go/src/rakuten.co.jp
ARG PROJECT

WORKDIR /usr/local/bin
COPY --from=build-env ${GOPREFIX}/${PROJECT} .
ENTRYPOINT [ "./app" ]
