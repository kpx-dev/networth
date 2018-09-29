# networth.app - Schema / Architecture Design

## DynamoDB Table Schema

Table: `networth`

All tokens for a user

```json
{
  "id": "demo@networth.app:token",
  "sort": "all",
  "tokens": [
    {
      "item_id": "pVVqaPm8MxFzn5MVk9JBIjq1WyDxM9FL1j3gQ",
      "access_token": "AQICAHs=...",
      "institution_name": "Bank of America",
      "institution_id": "ins_1",
    },
    {
      "item_id": "xxxxx",
      "access_token": "AQICAHifbIs=...",
      "institution_name": "Chase",
      "institution_id": "ins_2",
    }
  ],
}
```

Tokens for a specific institution

```json
{
  "id": "demo@networth.app:token",
  "sort": "ins_1",
  "tokens": [
    {
      "item_id": "pVVqaPm8MxFzn5MVk9JBIjq1WyDxM9FL1j3gQ",
      "access_token": "AQICAHs=...",
      "institution_name": "Bank of America",
      "institution_id": "ins_1",
    },
    {
      "item_id": "xxxxx",
      "access_token": "AQICAHifbIs=...",
      "institution_name": "Bank of America",
      "institution_id": "ins_1",
    }
  ],
}
```

Accounts for a user

```json
{
  "id": "demo@networth.app:account",
  "sort": "all",
  "accounts": [
    {
      "id": "xxx",
      "balances": {
        "available": 1234,
        "current": 123,
        "limit": 12,
        "iso_currency_code": "usd",
        "unofficial_currency_code": "usd",
      },
      "mask": "x123",
      "name": "Checking Core",
      "official_name": "Checking Core",
      "subtype": "",
      "type": "",
      "institution_id": "ins_1",
    },
  ],
}
```

Accounts for a user at a bank

```json
{
  "id": "demo@networth.app:account",
  "sort": "ins_1",
  "accounts": [
    {
      "id": "xxx",
      "balances": {
        "available": 1234,
        "current": 123,
        "limit": 12,
        "iso_currency_code": "usd",
        "unofficial_currency_code": "usd",
      },
      "mask": "x123",
      "name": "Checking Core",
      "official_name": "Checking Core",
      "subtype": "",
      "type": "",
      "institution_id": "ins_1",
    },
  ],
}
```

Net Worth for a user

```json
{
  "id": "demo@networth.app:networth",
  "sort": "2018-01-01T02:03:03Z",
  "networth": 150123,
}
```
