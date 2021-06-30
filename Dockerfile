FROM golang:1.16-alpine AS build

RUN apk add --no-cache git

WORKDIR /src/
COPY go.* .

RUN go mod download

COPY . .

# COPY cmd/server/main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/server cmd/server/main.go
RUN CGO_ENABLED=0 go build -o /bin/migration cmd/migration/main.go
RUN CGO_ENABLED=0 go build -o /bin/local-dev-setup cmd/local_dev_setup/main.go

FROM scratch
COPY --from=build /bin/server bin/server
COPY --from=build /bin/migration bin/migration
COPY --from=build /bin/local-dev-setup bin/local-dev-setup
COPY .env .
ENTRYPOINT ["/bin/akc"]
