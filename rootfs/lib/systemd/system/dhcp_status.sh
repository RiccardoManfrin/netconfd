#!/bin/sh
#Yocto busybox udhcpc systemd based DHCP service start/stop

ETH=$1
export $(systemctl show --property=ActiveState udhcpc@$ETH)

echo -en $ActiveState
