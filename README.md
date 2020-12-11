OpenAPI templates
-----------------

## Developing

Entrypoint to increment OpenAPI spec is 

    swaggerui/openapi.yaml

Extend it and run `./gen_templates.sh` to regenerate templates. Than manually integrate the relevant parts and discartd the rest :(

Pushing on git will create an artifact with a docker image embedding the binary.

Build locally with `make`

Code coverage
-------------

![coverage](https://gitlab.com/gitlab-org/gitlab/badges/master/coverage.svg?job=coverage)

ToDo list
---------

## DHCP static addr Warning emit

## MTU set

This is needed by UPF

## Golden config

Need to translate the config from networkd into mine from here

https://gitlab.lan.athonet.com:8443/core/meta-athonet/tree/master/recipes-core/systemd/files

## Use DBUS not scripts

https://gitlab.lan.athonet.com:8443/core/ncm/blob/master/n1dbus/n1dbus.go

https://github.com/godbus/dbus/

## Eth1 GRO e GSO

Currently networkd writes this file: https://gitlab.lan.athonet.com:8443/core/meta-athonet/blob/master/recipes-core/systemd/files/20-eth1.link

	TCPSegmentationOffload=false 
	TCP6SegmentationOffload=false 
	GenericSegmentationOffload=false 
	GenericReceiveOffload=false 
	LargeReceiveOffload=false

Since this appears as a very specific need of the current UPF after a chat with Davide we reasoned that it is probably 
conevinent to keep this enforcement in the UPF recipe itself for now.

## Run containerized 

This is the egg and chicken problem (podman and network bootstrap)

## Check tun device creation

Tun device tun0 is used by the UPF. I can create it as courtesy, although UPF does create it if it does not find it.

## Unit tests automation

Spin yocto/alpine machine dedicated to unit tests.

## non-local-bind ?

https://www.cyberciti.biz/faq/linux-bind-ip-that-doesnt-exist-with-net-ipv4-ip_nonlocal_bind/
https://www.kernel.org/doc/Documentation/networking/ip-sysctl.txt
http://web.mit.edu/rhel-doc/4/RH-DOCS/rhel-rg-it-4/s1-proc-sysctl.html

## CNI
https://www.nuagenetworks.net/blog/container-networking-standards/
https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/
https://github.com/containers/dnsname/blob/master/README_PODMAN.md
https://github.com/containers/dnsname/blob/master/plugins/meta/dnsname/main.go
https://github.com/containernetworking/plugins/blob/master/plugins/sample/main.go
