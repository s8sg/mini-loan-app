FROM golang:1.20

WORKDIR /go/src/github.com/s8sg/mini-loan-app/app

COPY go.mod         go.mod
COPY go.sum         go.sum
COPY main.go        .

RUN CGO_ENABLED=1  go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/s8sg/mini-loan-app/app ./
CMD ["./app"]