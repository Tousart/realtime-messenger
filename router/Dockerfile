FROM golang:1.25-alpine AS build

WORKDIR /build

COPY ./go.mod ./go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine:latest AS app

WORKDIR /app

COPY --from=build /build/main ./main

ENTRYPOINT ["/app/main"]