## Build Base
FROM golang:1.23-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

## Build api binary
FROM base AS api-binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -ldflags="-s -w" ./cmd/api
RUN chmod +x server

## Build api
FROM scratch AS api
COPY --from=api-binary /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]

## Build backend binary
FROM base AS backend-binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -ldflags="-s -w" ./cmd/backend
RUN chmod +x server

## Build backend
FROM scratch AS backend
COPY --from=backend-binary /app/server /server
EXPOSE 8081
ENTRYPOINT ["/server"]
