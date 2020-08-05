FROM golang:1-alpine as builder
RUN apk --no-cache add build-base git gcc
WORKDIR /go/src/github.com/jasonmccallister/httputil
COPY . .
RUN go build -o httputil

FROM alpine:3.12
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/jasonmccallister/httputil/httputil .
ENTRYPOINT ./httputil
EXPOSE 8000
