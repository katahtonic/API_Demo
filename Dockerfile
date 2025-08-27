FROM golang:1.23

WORKDIR /app

COPY . .

RUN go build -x -o main .

EXPOSE 8080

CMD ["./main"]