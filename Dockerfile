FROM golang:1.23 AS build-stage

# set destination for COPY
WORKDIR /app

# download Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

FROM scratch AS build-release-stage

WORKDIR /app

COPY --from=build-stage /api /api

EXPOSE 8080

CMD ["/api"]

