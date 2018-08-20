#!/usr/bin/env bash

aws cloudformation create-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml
