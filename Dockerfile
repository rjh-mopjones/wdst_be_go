FROM golang:1.22

WORKDIR /go/src/app
COPY . .

RUN go mod tidy
RUN go build
ENTRYPOINT ["./wdst_be"]