FROM golang:latest

WORKDIR /data
RUN git clone https://github.com/doublepi123/mmssbackend
WORKDIR /data/mmssbackend
RUN go build .

EXPOSE 58888
ENTRYPOINT ["/data/mmssbackend/mmssbackend"]
