module gitlab.lan.athonet.com/core/netconfd

go 1.16

replace github.com/getkin/kin-openapi => github.com/riccardomanfrin/kin-openapi v0.22.2-0.20210421143010-b343c1d8f9c9

replace github.com/vishvananda/netlink => github.com/riccardomanfrin/netlink v1.1.1-0.20210416214702-2421371a9da7

require (
	github.com/getkin/kin-openapi v0.53.0
	github.com/gorilla/mux v1.8.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5
	golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57
)
