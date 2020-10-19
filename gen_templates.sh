#!/bin/bash
#DOWNLOAD TEMPLATES:
#templates/go-server taken from 
#https://github.com/OpenAPITools/openapi-generator/tree/master/modules/openapi-generator/src/main/resources
#GENERAGE TEMPLATES
#docker run  \
#    -v "${PWD}:/local" \
#    -u $(id -u):$(id -g) \
#    --rm \
#    openapitools/openapi-generator-cli author template \
#    -g go-server \
#    -o /local/templates/go-server/

rm server -rf
mkdir server


docker run \
    -e GO_POST_PROCESS_FILE="/usr/local/bin/gofmt -w" \
    -u $(id -u):$(id -g) \
    --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli generate \
    -i /local/swaggerui/swagger.yaml -t /local/templates/go-server/ \
    -g go-server -o /local/server 2>&1 > /dev/null

docker run \
    -e GO_POST_PROCESS_FILE="/usr/local/bin/gofmt -w" \
    -u $(id -u):$(id -g) \
    --rm \
    -v "${PWD}:/local" \
    openapitools/openapi-generator-cli generate \
    -i /local/swaggerui/swagger.yaml \
    -g go -o /local/server/go-client/ 2>&1 > /dev/null

cp server/go-client/model_*.go server/go/
cp server/go-client/utils.go server/go/
rm server/go-client -rf