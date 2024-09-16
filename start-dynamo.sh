#!/usr/bin/env bash


echo starting container
docker run --rm -d -p  8000:8000 amazon/dynamodb-local

echo creating table
aws dynamodb create-table \
  --endpoint-url http://localhost:8000 \
  --table-name Books \
  --attribute-definitions AttributeName=ISBN,AttributeType=S \
  --key-schema AttributeName=ISBN,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST


