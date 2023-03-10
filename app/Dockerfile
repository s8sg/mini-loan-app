FROM golang:1.20

WORKDIR /go/src/github.com/s8sg/mini-loan-app/app

COPY . .

RUN CGO_ENABLED=1  go build -a -installsuffix cgo -o loanapp .

CMD ["./loanapp"]