# networth.app - Schema / Architecture Design

## DynamoDB Table Schema

Table: `networth`

All tokens for a user

```json
{
  "key": "demo@networth.app:token",
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
  "key": "demo@networth.app:token",
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
  "key": "demo@networth.app:account",
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
  "key": "demo@networth.app:account",
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














<!--
Table: `networth_token`
| email (partition key: string) | tokens (map) |
| --- | --- | --- |
| demo@networth.app |
```json
{ "ins_1":
  {
    "name": "Bank of America",
    "tokens": ["encrypted_token_1", "encrypted_token_2"],
    "accounts": [...]
  }
}
```

Table: `networth_history`

| email (partition key: string) | datetime (sort key: string) | networth (number) |
| --- | --- | --- |
| demo@networth.app | 2018-01-01T12:01:02 | 100 |
| demo@networth.app | 2018-02-01T12:01:02 | 140 |
| demo@networth.app | 2018-02-02T12:01:02 | 240 |

Table: `networth_transaction`
| email (partition key: string) | datetime (sort key: string) | name | category | value
| --- | --- | --- | --- | --- |
| demo@networth.app | 2018-01-01T12:01:02 | McDonald purchase | Food | $10
| demo@networth.app | 2018-01-01T12:02:02 | networth.app purchase | SASS | $100 -->
