#!/bin/bash
rm server -rf
mkdir server
#templates/go-server taken from 
#https://github.com/OpenAPITools/openapi-generator/tree/master/modules/openapi-generator/src/main/resources

docker run \
    -e GO_POST_PROCESS_FILE="/usr/local/bin/gofmt -w" \
    -u $(id -u):$(id -g) \
    --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli generate \
    -i /local/swaggerui/swagger.yaml -t /local/templates/go-server/ \
    -g go-server -o /local/server 2>&1 > /dev/null