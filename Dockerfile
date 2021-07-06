FROM golang:1.16-buster as build

WORKDIR /src/
COPY go.* .

RUN go mod download

COPY . .

# COPY cmd/server/main.go go.* /src/
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /bin/server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /bin/migration cmd/migration/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /bin/local-dev-setup cmd/local_dev_setup/main.go

FROM debian:buster-slim
COPY --from=build /bin/server bin/server
COPY --from=build /bin/migration bin/migration
COPY --from=build /bin/local-dev-setup bin/local-dev-setup
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY .env .
ENTRYPOINT ["/bin/server"]
