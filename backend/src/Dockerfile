FROM golang:1.21.5

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src /app/src

WORKDIR /app/src/cmd

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
