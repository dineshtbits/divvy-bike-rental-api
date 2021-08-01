# Divvy Bike Rental API

## How to run this?

```
go run cmd/main.go 
```
It shows all available endpoints. Here are some sample curls

Get station information
```
curl --location --request GET 'http://localhost:8081/api/v1/stations/4' \
--header 'Authorization: Basic YWRtaW46YWRtaW4='
```
Get Rider summary
```
curl --location --request POST 'http://localhost:8081/api/v1/trips/riders/summary' \
--header 'Authorization: Basic YWRtaW46YWRtaW4=' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "filters": {
        "station_ids": [56, 59, 174]
    }
}'
```
Trip Summary
```
curl --location --request POST 'http://localhost:8081/api/v1/trips/summary' \
--header 'Authorization: Basic YWRtaW46YWRtaW4=' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "filters": {
        "station_ids": [56, 59, 174]
    }
}'
```

## How to run tests?

```
cd cmd && go test
```