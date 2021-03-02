FROM golang:alpine3.12
WORKDIR /go/src/github.com/abhide/envoy-access-log-client/
COPY main.go .
COPY go.mod .
COPY go.sum .
RUN go build -o envoy-access-log-client ./main.go

FROM alpine:3.12
WORKDIR /root/
COPY --from=0 /go/src/github.com/abhide/envoy-access-log-client/envoy-access-log-client .
CMD ["./envoy-access-log-client"]
