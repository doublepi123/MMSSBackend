FROM golang:latest

WORKDIR /data
RUN git clone https://github.com/doublepi123/MMSSBackend
WORKDIR /data/MMSSBackend
RUN go build .

EXPOSE 58888
ENTRYPOINT ["./MMSSBackend"]
