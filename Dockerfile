FROM golang:latest

WORKDIR /data
RUN git clone https://github.com/doublepi123/MMSSbackend
WORKDIR /data/MMSSbackend
RUN go build .

EXPOSE 58888
ENTRYPOINT ["./MMSSbackend"]
