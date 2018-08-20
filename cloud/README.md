# networth.app Cloud Setup

## AWS

```shell
# Launching AWS infrastructure
./bin/create-aws-infra.sh

# Update AWS infrastructure
./bin/update-aws-infra.sh
```

## DNS

MX records:

```shell
5 gmr-smtp-in.l.google.com.
10 alt1.gmr-smtp-in.l.google.com.
20 alt2.gmr-smtp-in.l.google.com.
30 alt3.gmr-smtp-in.l.google.com.
40 alt4.gmr-smtp-in.l.google.com.
```
