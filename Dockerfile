ARG GO_VERSION=1.24

FROM golang:$(GO_VERSION) AS build-stage
# set destination for COPY
WORKDIR /app

# download Go modules and dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM scratch AS build-release-stage

COPY --from=build-stage /bin/api /bin/api

EXPOSE 80

CMD ["/bin/api"]

