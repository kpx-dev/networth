#!/usr/bin/env bash

aws ssm put-parameter --type String --overwrite --name $1 --value $2
