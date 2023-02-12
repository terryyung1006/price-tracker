# BTCUSD Price Tracker

This price tracker will get USD price of BTC per minute from Coingecko public API. Exchange rate will be stored in mongoDB every 5 minutes. The latest exchange rate, exchange rate at specific timestamp and average price within time range can be get by specific API interface provided below.

localhost:8081 will be the mongo-express panel to view the stored records. Records can be shown after 5 mins(new records stored in DB every 5 mins).

if running server without docker needed, please comment out price-tracker service in compose.yml, compose up with compose.yml, and then run command

```bash
  go run main.go
```


## Running Tests

To run tests, run the following command in repo folder

```bash
  go mod tidy
  go test ./test/api_test.go -v
```


## Run Locally

Start the server

```bash
  docker-compose -f compose.yml up
```

# RESTAPIDocs

API document for price related open endpoints

## Open Endpoints

Open endpoints require no Authentication.

* [LatestPrice](getLatestPrice.md) : `GET localhost:8080/api/last_price/`

* [PriceByTimestamp](getPriceByTimestamp.md) : `GET localhost:8080/api/price_by_timestamp/`

* [AveragePrice](getAveragePriceInRange.md) : `GET localhost:8080/api/average_price_in_range/`

# Hierarchy

## Worker

Cronjob defined simple cronjob runner.

Workers in worker folder implement cronjob with different propose.
Retryer workers will handle failed cases in normal workers.

## API Handlers

API interface implemented with controller, service and repository.

## http

Connection with other public API(Coingecko) handled in lib/http folder.

## test

API handlers test cases in test folder, used testify.

## variable

change cronjob interval in variable folder. DB storing interval will also be updated according. just a simple way to handle global variable.
