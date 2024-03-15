FROM docker.io/golang:alpine as builder
WORKDIR /root
ADD . .
RUN go mod init ship && go mod tidy
RUN env CGO_ENABLED=0 go build -v -ldflags '-w -s'
FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/ship .
ENTRYPOINT ["/root/ship"]
