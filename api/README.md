# networth.app - REST API

## Install

This api requires Golang 1.11 and up

```shell
# install packages
go mod download
```

## Start service

```shell
# install gin
go get -u github.com/gin-gonic/gin

# start service locally
make start-api

# make a healthcheck request using httpie
http :3000/healthcheck
```

## Testing

Docker is required for test as it will launch a local-dynamodb instance in a container.

```shell
# pull dynamodb-local
docker pull amazon/dynamodb-local

# start a dynamodb-local container
docker run --rm -d -p 8000:8000 amazon/dynamodb-local

# test api
make test
```

## Deploy service

```shell
make deploy-api
```
