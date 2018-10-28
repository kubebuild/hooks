FROM golang:1.11 AS builder
WORKDIR /go/src/github.com/kubebuild/webhooks/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/kubebuild/webhooks/app .
ENTRYPOINT ["./app"]