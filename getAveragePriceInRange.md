# Get Average Price In Range

Get average price in given timestamp range.

**URL** : `/api/price_by_timestamp/`

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

**Data example**

```json
{
    "pair_tag": "btcusd",
    "timestamp_from": "1676179815",
    "timestamp_to": "1676180295"
}
```

## Success Response

**Code** : `200 OK`

**Content example**

```json
{
    "data": "21000"
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