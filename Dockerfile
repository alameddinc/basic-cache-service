FROM golang:1.15.6

RUN mkdir /server

ADD . /server

WORKDIR /server

RUN go build -o main app/main.go

CMD ["main"]