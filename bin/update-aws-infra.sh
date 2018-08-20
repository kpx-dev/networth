#!/usr/bin/env bash

aws cloudformation update-stack --stack-name networth --capabilities CAPABILITY_IAM --template-body file://cloud/aws.infra.yml
