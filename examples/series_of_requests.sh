#!/usr/bin/env sh
set -e

export YAHR_CONFIG_FILE="./httpbin.yaml"

resp=$(yahr run -s httpbin post)

id=$(jq .json.id)
if [ -z "${id}" ]; then
    echo "Failed to get id from newly created object"
    exit 1
fi

echo "Created object ${id}"

resp=$(OBJECT_ID=${id} yahr run -s httpbin put)
yahr=$(jq .json.yahr)

echo "Successfully put object ${id} with ${yahr}"

OBJECT_ID=${id} yahr run -s httpbin delete

echo "Successfully deleted object ${id}"
