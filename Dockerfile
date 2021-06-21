FROM golang:1.16-alpine AS build

RUN apk add --no-cache git

WORKDIR /src/
COPY go.* .

RUN go mod download

COPY . .

# COPY cmd/server/main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/akc cmd/server/main.go

FROM scratch
COPY --from=build /bin/akc bin/akc
ENTRYPOINT ["/bin/akc"]
