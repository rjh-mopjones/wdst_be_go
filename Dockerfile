FROM golang:1.22

WORKDIR /go/src/app
COPY . .

RUN go build -o wdst_be
RUN chmod +x wdst_be


ENTRYPOINT ["./wdst_be"]