#!/bin/bash
rm server -rf

#templates/go-server taken from 
#https://github.com/OpenAPITools/openapi-generator/tree/master/modules/openapi-generator/src/main/resources
docker run -u $(id -u):$(id -g) --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
    -i /local/swaggerui/swagger.yaml -t /local/templates/go-server/ \
    -g go-server -o /local/server