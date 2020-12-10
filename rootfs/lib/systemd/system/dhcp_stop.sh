#!/bin/sh
#Yocto busybox udhcpc systemd based DHCP service start/stop

ETH=$1
export $(systemctl show --property=ActiveState udhcpc@$ETH)

if [ $ActiveState == "active" ]; then
    systemctl stop udhcpc@$ETH
    export $(systemctl show --property=ActiveState udhcpc@$ETH)
    echo -en "Service stopped; state: $ActiveState"
else
    echo -en "Service not running"
fi