#!/usr/bin/env bash

aws cloudformation validate-template --template-body file://$1
