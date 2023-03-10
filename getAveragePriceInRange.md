# Get Average Price In Range

Get average price in given timestamp range.

**URL** : `localhost:8080/api/price_by_timestamp/`

**Method** : `GET`

**Auth required** : NO

**Query Parameters**

```json
{
    "pair_tag": "[tag string for pair, only accept btcusd now]",
    "timestamp_from": "[average price from this given timestamp]",
    "timestamp_to": "[average price to this given timestamp]",
}
```

**Example**

```bash
    localhost:8080/api/average_price_in_range?timestamp_from=1676220328&timestamp_to=1676220568&pair_tag=btcusd
```

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
    "data": "22115.955"
}
```

## Error Response

**Condition** : If database doesn't have enough data to satisfy given time range.

**Code** : `416`

**Content** :

```json
{
    "message": "input range not satisfiable"
}
```