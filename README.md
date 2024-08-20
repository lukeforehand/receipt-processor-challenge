# Receipt Processor

## Instructions

### Running the application

```bash
  docker compose up --build --remove-orphans
```

### Unit Testing

```bash
  go test ./...
```

## Notes

- generated server with `oapi-codegen`
- `docker compose` to ease local development
- docker image size is minimized to the go static binary
- `/health` endpoint for liveness probes
- assuming SSL termination at the load balancer
- assuming an authentication proxy so no auth middleware

## Load Testing Results

Install pocket load tester

```bash
  go install github.com/vearutop/plt@latest
```

### POST /receipts/process

```bash
  
  plt --live-ui --duration=20s --rate-limit=60 curl -X POST "http://localhost:8080/receipts/process" -d'{
    "retailer": "Target",
    "purchaseDate": "2022-01-02",
    "purchaseTime": "13:13",
    "total": "1.25",
    "items": [
        {"shortDescription": "Pepsi - 12-oz", "price": "1.25"}
    ]
  }'

  Requests per second: 60.05
  Successful requests: 1201
  Time spent: 20.002s

  Request latency percentiles:
  99%: 1.46ms
  95%: 0.81ms
  90%: 0.71ms
  50%: 0.52ms

  Response samples (first by status code):
  [HTTP/1.1 200]
  Content-Length: 46
  Content-Type: application/json
  Date: Tue, 20 Aug 2024 05:11:44 GMT

  {"id":"7d4d837b-ef5e-47c0-89a9-889657b66eb9"}

```

### GET /receipts/{id}/points

```bash
  plt --live-ui --duration=20s --rate-limit=60 curl -X GET "http://localhost:8080/receipts/7d4d837b-ef5e-47c0-89a9-889657b66eb9/points"

  Requests per second: 60.04
  Successful requests: 1201
  Time spent: 20.002s

  Request latency percentiles:
  99%: 0.86ms
  95%: 0.67ms
  90%: 0.63ms
  50%: 0.45ms

  Response samples (first by status code):
  [HTTP/1.1 200]
  Content-Length: 14
  Content-Type: application/json
  Date: Tue, 20 Aug 2024 05:12:20 GMT

  {"points":31}
```
