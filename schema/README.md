# networth.app - Schema / Architecture Design

## DynamoDB Table Schema

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
| demo@networth.app | 2018-01-01T12:02:02 | networth.app purchase | SASS | $100
