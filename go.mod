module github.com/riccardomanfrin/netconfd

// +heroku goVersion go1.16
go 1.16

replace github.com/getkin/kin-openapi => github.com/riccardomanfrin/kin-openapi v0.22.2-0.20210421143010-b343c1d8f9c9

replace github.com/vishvananda/netlink => github.com/riccardomanfrin/netlink v1.1.1-0.20210416214702-2421371a9da7

require (
	github.com/getkin/kin-openapi v0.61.0
	github.com/gorilla/mux v1.8.0
	github.com/heroku/x v0.0.29
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/vishvananda/netlink v1.1.0
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015
)
