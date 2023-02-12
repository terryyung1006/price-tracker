# BTCUSD Price Tracker

This price tracker will get USD price of BTC per minute from Coingecko public API. Exchange rate will be stored in mongoDB every 5 minutes. The latest exchange rate, exchange rate at specific timestamp and average price within time range can be get by specific API interface provided below.




## Running Tests

To run tests, run the following command in repo folder

```bash
  go test ./test/api_test.go -v
```


## Run Locally

Start the server

```bash
  sudo docker-compose -f compose.yml up
```

# RESTAPIDocs

API document for price related open endpoints

## Open Endpoints

Open endpoints require no Authentication.

* [LatestPrice](getLatestPrice.md) : `POST /api/last_price/`

* [PriceByTimestamp](getPriceByTimestamp.md) : `POST /api/price_by_timestamp/`

* [AveragePrice](getAveragePriceInRange.md) : `POST /api/average_price_in_range/`