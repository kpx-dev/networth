# networth.app - REST API

## Install

This api requires Golang 1.11 and up

```shell
# install packages
go mod download
```

## Start service

```shell
# start service locally
make start-api

# make a healthcheck request using httpie
http :3000/healthcheck
```

## Deploy service

```shell
make deploy-api
```
