FROM golang:1.18-alpine3.14

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build main.go

EXPOSE 3000

CMD ["./main"]