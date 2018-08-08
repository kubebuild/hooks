FROM golang:latest
WORKDIR /go/src/github.com/kubebuild/webhooks/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/kubebuild/webhooks/app .
CMD ["./app"]  