# Get Price By Timestamp

Get price at timestamp nearest to given value. Will get from Coingecko public API if no record in Database.

**URL** : `localhost:8080/api/price_by_timestamp/`

**Method** : `GET`

**Auth required** : NO

**Query Parameters**

```json
{
    "pair_tag": "[tag string for pair, only accept btcusd now]",
    "timestamp": "[price at timestamp]",
}
```

**Example**

```bash
    localhost:8080/api/price_by_timestamp?timestamp=1676118169&pair_tag=btcusd
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

**Condition** : If given timestamp is not valid.

**Code** : `400`

**Content** :

```json
{
    "message": "timestamp [1676179815] length of digit invalid"
}
```