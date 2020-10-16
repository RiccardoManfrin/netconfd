#!/bin/bash

docker network create --driver bridge nc_net
docker run -v $(pwd):/app --rm -it --name nc_sandbox --network nc_net --cap-add NET_ADMIN alpine /bin/sh
#docker run -v $(pwd):/app --rm -dit --name nc_sandbox --network nc_net --cap-add NET_ADMIN alpine /bin/sh -c "/app/netconfd -log - -config /app/netconfd.json"
echo "Start testing [press ENTER to unstage test env]"
read DONE
docker rm --force nc_sandbox
docker network rm nc_net