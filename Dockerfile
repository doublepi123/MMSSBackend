FROM golang:latest

WORKDIR $GOPATH/src/github.com/doublepi123/mmssbackend
COPY . $GOPATH/src/github.com/doubelpi123/mmssbackend
RUN go build .

EXPOSE 58888
ENTRYPOINT ["./mmssbackend"]
