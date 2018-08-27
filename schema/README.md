# networth.app - Schema / Architecture Design

## DynamoDB Table Schema

Table: `networth_token`
| email (partition key: string) | tokens (list) | accounts (list) |
| --- | --- | --- |
| demo@networth.com | [ sandbox-access-token-abc, sandbox-access-token-xyz ] | [ account-01, account-02 ] |
| user2@networth.com | [ sandbox-access-token-zzz ] | [ boa-01, chase-02 ] |

Table: `networth_history`

| email (partition key: string) | datetime (sort key: string) | networth (number) |
| --- | --- | --- |
| demo@networth.com | 2018-01-01T12:01:02 | 100 |
| demo@networth.com | 2018-02-01T12:01:02 | 140 |
| demo@networth.com | 2018-02-02T12:01:02 | 240 |

Table: `networth_transaction`
| email (partition key: string) | datetime (sort key: string) | name | category | value
| --- | --- | --- | --- | --- |
| demo@networth.com | 2018-01-01T12:01:02 | McDonald purchase | Food | $10
| demo@networth.com | 2018-01-01T12:02:02 | networth.app purchase | SASS | $100
