#!/usr/bin/env bash

aws s3 cp web/index.html s3://demo.networth.app/
aws s3 cp web/assets s3://demo.networth.app/ --recursive
