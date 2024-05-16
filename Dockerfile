FROM golang:1.22

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /wdst_be

RUN chmod +x /wdst_be


ENTRYPOINT ["/wdst_be"]
#ENTRYPOINT ["tail", "-f", "/dev/null"]