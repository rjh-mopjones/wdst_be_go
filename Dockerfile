FROM golang:1.22

WORKDIR /go/src/app
COPY . .

RUN go build
CMD ["./wdst_be"]
