FROM golang:alpine

RUN apk add git

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main /app/src/main.go

EXPOSE 8080
CMD ["/app/main"]
