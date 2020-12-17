# OpenAPI templates

## Developing

Entrypoint to increment OpenAPI spec is 

    swaggerui/openapi.yaml

Extend it and run `./gen_templates.sh` to regenerate templates. Than manually integrate the relevant parts and discartd the rest :(

Pushing on git will create an artifact with a docker image embedding the binary.

Build locally with `make`

# Code coverage

![coverage](https://gitlab.com/gitlab-org/gitlab/badges/master/coverage.svg?job=coverage)

# R1 ToDo list

## DHCP static addr Warning emit [2h]

## MTU set [3h]

This is needed by UPF

## Golden config [3h]

Need to translate the config from networkd into mine from here

https://gitlab.lan.athonet.com:8443/core/meta-athonet/tree/master/recipes-core/systemd/files

## Use DBUS not scripts [1d - risk]

Library: https://github.com/godbus/dbus/

Our code: https://gitlab.lan.athonet.com:8443/core/ncm/blob/master/n1dbus/n1dbus.go

Autogen constants command

    busctl introspect org.freedesktop.systemd1 /org/freedesktop/systemd1

Risk assessment: unknown technology

## Eth1 GRO e GSO [1d if it goes smoothly - mandatory - unknown requirement]

Currently networkd writes this file: https://gitlab.lan.athonet.com:8443/core/meta-athonet/blob/master/recipes-core/systemd/files/20-eth1.link

	TCPSegmentationOffload=false 
	TCP6SegmentationOffload=false 
	GenericSegmentationOffload=false 
	GenericReceiveOffload=false 
	LargeReceiveOffload=false

This appears as a very specific need of the current UPF after a chat with Davide: UPF needs it otherwise packets are lost in the kernel. Proposal was to confine the enforcement of the configuration within the UPF. Carlo disagrees. Needs decision

## Check tun device creation [2h]

Tun device tun0 is used by the UPF. I can create it as courtesy, although UPF does create it if it does not find it.

## Documentation [2 or 3 days]

## Unit tests automation [2d]

Spin yocto/alpine machine dedicated to unit tests.

## Run containerized [Needs investigation]

This is the egg and chicken problem (podman and network bootstrap)

## Yocto integration [1d - risks]

Needs to write the bb recipe for yocto [containerized or not makes a difference here]

### Enablement of VPN

    sed -i -e 's/dev\ tun/dev\ tunvpn0/g' /etc/openvpn/ovpn.conf
    echo -en "ipchange '/bin/sh -c \"ip link del tunvpn0\"'" >> /etc/openvpn/ovpn.conf

Risk assessment: 1d of work but if it takes 1 day to build I can have multiple iterations to converge

# Long Term Investigation

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
