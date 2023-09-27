FROM golang:1.21

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o user-microservice ./cmd/main

EXPOSE 80

ENTRYPOINT [ "./user-microservice" ]