## Build Stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# Download module dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build static binary
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server -ldflags="-s -w"
RUN chmod +x server

## Create a tiny image
FROM scratch

COPY --from=build /app/server /server
COPY --from=build /app/api.yml /api.yml
COPY --from=build /app/static /static
EXPOSE 8080
ENTRYPOINT ["/server"]
