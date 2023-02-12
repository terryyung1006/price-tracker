# Get Latest Price

Get latest price get by worker from Coingecko public api.

**URL** : `/api/last_price/`

**Method** : `GET`

**Auth required** : NO

**Query Parameters**

```json
{
    "pair_tag": "[tag string for pair, only accept btcusd now]"
}
```

**Data example**

```json
{
    "pair_tag": "btcusd",
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

**Condition** : If 'price_tag' is missing.

**Code** : `400 BAD REQUEST`

**Content** :

```json
{
    "message": "only btcusd is supported now"
}
```