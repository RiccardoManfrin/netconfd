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

## Move code into nc

    func (s *NetworkApiService) ConfigLinkCreate
	

## Use DBUS not scripts [1d - risk]

Library: https://github.com/godbus/dbus/

Our code: https://gitlab.lan.athonet.com:8443/core/ncm/blob/master/n1dbus/n1dbus.go

Autogen constants command

    busctl introspect org.freedesktop.systemd1 /org/freedesktop/systemd1

Risk assessment: unknown technology

## Documentation [2 or 3 days]

## Unit tests automation [2d]

Spin yocto/alpine machine dedicated to unit tests.

## Run containerized [Needs investigation]

This is the egg and chicken problem (podman and network bootstrap)

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
