#!/bin/bash


sudo ip netns add nc_net
sudo ip link add mymacvlan1 link wlan0 type macvlan mode bridge
sudo ifconfig mymacvlan1 192.168.178.79 up
sudo ip link set mymacvlan1 netns nc_net


sudo podman run -v $(pwd):/app --rm -it --name nc_sandbox --network ns:/var/run/netns/nc_net --cap-add NET_ADMIN alpine /bin/sh

sudo ifconfig mymacvlan1 192.168.178.79 up

#docker run -v $(pwd):/app --rm -dit --name nc_sandbox --network nc_net --cap-add NET_ADMIN alpine /bin/sh -c "/app/netconfd -log - -config /app/netconfd.json"
echo "Start testing [press ENTER to unstage test env]"
read DONE

#docker rm --force nc_sandbox
#sudo ip link del mymacvlan1
sudo ip netns del nc_net