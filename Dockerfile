FROM golang:1.23 AS builder

WORKDIR /app

# copy Go binary files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy all Go files into the container
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/api ./cmd/api \
    && go build -o ./bin/migrates ./cmd/migrates

EXPOSE 8080

CMD [ "/app/bin/api" ]

