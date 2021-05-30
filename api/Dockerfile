FROM golang:1.14.4-alpine3.12 AS builder
WORKDIR /go/src/github.com/manufacturer/api
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/api -mod vendor github.com/manufacturer/api/cmd

FROM alpine:3.12
COPY --from=builder /bin/api /opt

ENV TZ Europe/Moscow
RUN apk add --no-cache tzdata \
    && ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

CMD ["/opt/api"]
