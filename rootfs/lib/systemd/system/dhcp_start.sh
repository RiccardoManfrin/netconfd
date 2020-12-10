#!/bin/sh
#Yocto busybox udhcpc systemd based DHCP service start/stop

ETH=$1
export $(systemctl show --property=ActiveState udhcpc@$ETH)

if [ $ActiveState == "inactive" ]; then
    systemctl start udhcpc@$ETH
    export $(systemctl show --property=ActiveState udhcpc@$ETH)
    echo "Service started; state: $ActiveState"
else
    echo "Service running already"
fi