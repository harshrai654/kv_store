# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as build
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/main ./cmd/web
COPY .env ./build

FROM scratch
COPY --from=build /app/build/ .
ENTRYPOINT [ "./main" ]




