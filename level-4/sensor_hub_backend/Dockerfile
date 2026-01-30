FROM golang:alpine3.23

WORKDIR /go/src/app

COPY . .

RUN go build -o main main.go

CMD ["./main"]

EXPOSE 8080