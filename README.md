# app-metadata-server

To run the server:
```
go run main.go
```

Insert an app:
```
curl -X POST http://localhost:8000/apps --data-binary '@payloads/valid1.yaml'
```

Searching apps:
```
# Retrieve all apps
curl http://localhost:8000/apps

# Filter by title (exact match)
curl "http://localhost:8000/apps?title=Valid%20App%201"

# Filter by title AND version
curl "http://localhost:8000/apps?title=Valid%20App%201&version=0.0.1"

# Filter by description substring match
curl "http://localhost:8000/apps?descriptionContains=Interesting"
```

To run the test suite:
```
go test ./server
```

To run benchmarks:
```
go test ./server --bench=.
```
