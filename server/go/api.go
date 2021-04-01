/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.3.0
 * Contact: support@athonet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"net/http"
)

// NetworkApiRouter defines the required methods for binding the api requests to a responses for the NetworkApi
// The NetworkApiRouter implementation should parse necessary information from the http request,
// pass the data to a NetworkApiServicer to perform the required actions, then write the service results to the http response.
type NetworkApiRouter interface {
	ConfigDHCPCreate(http.ResponseWriter, *http.Request)
	ConfigDHCPDel(http.ResponseWriter, *http.Request)
	ConfigDHCPGet(http.ResponseWriter, *http.Request)
	ConfigDHCPsGet(http.ResponseWriter, *http.Request)
	ConfigDNSCreate(http.ResponseWriter, *http.Request)
	ConfigDNSDel(http.ResponseWriter, *http.Request)
	ConfigDNSGet(http.ResponseWriter, *http.Request)
	ConfigDNSsGet(http.ResponseWriter, *http.Request)
	ConfigLinkCreate(http.ResponseWriter, *http.Request)
	ConfigLinkDel(http.ResponseWriter, *http.Request)
	ConfigLinkGet(http.ResponseWriter, *http.Request)
	ConfigLinksGet(http.ResponseWriter, *http.Request)
	ConfigRouteCreate(http.ResponseWriter, *http.Request)
	ConfigRouteDel(http.ResponseWriter, *http.Request)
	ConfigRouteGet(http.ResponseWriter, *http.Request)
	ConfigRoutesGet(http.ResponseWriter, *http.Request)
	ConfigUnmanagedCreate(http.ResponseWriter, *http.Request)
	ConfigUnmanagedDel(http.ResponseWriter, *http.Request)
	ConfigUnmanagedGet(http.ResponseWriter, *http.Request)
	ConfigUnmanagedListGet(http.ResponseWriter, *http.Request)
}
// SystemApiRouter defines the required methods for binding the api requests to a responses for the SystemApi
// The SystemApiRouter implementation should parse necessary information from the http request,
// pass the data to a SystemApiServicer to perform the required actions, then write the service results to the http response.
type SystemApiRouter interface {
	ConfigGet(http.ResponseWriter, *http.Request)
	ConfigPatch(http.ResponseWriter, *http.Request)
	ConfigSet(http.ResponseWriter, *http.Request)
	PersistConfig(http.ResponseWriter, *http.Request)
	ResetConfig(http.ResponseWriter, *http.Request)
}

// NetworkApiServicer defines the api actions for the NetworkApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type NetworkApiServicer interface {
	ConfigDHCPCreate(context.Context, Dhcp) (ImplResponse, error)
	ConfigDHCPDel(context.Context, string) (ImplResponse, error)
	ConfigDHCPGet(context.Context, string) (ImplResponse, error)
	ConfigDHCPsGet(context.Context) (ImplResponse, error)
	ConfigDNSCreate(context.Context, Dns) (ImplResponse, error)
	ConfigDNSDel(context.Context, Dnsid) (ImplResponse, error)
	ConfigDNSGet(context.Context, Dnsid) (ImplResponse, error)
	ConfigDNSsGet(context.Context) (ImplResponse, error)
	ConfigLinkCreate(context.Context, Link) (ImplResponse, error)
	ConfigLinkDel(context.Context, string) (ImplResponse, error)
	ConfigLinkGet(context.Context, string) (ImplResponse, error)
	ConfigLinksGet(context.Context) (ImplResponse, error)
	ConfigRouteCreate(context.Context, Route) (ImplResponse, error)
	ConfigRouteDel(context.Context, string) (ImplResponse, error)
	ConfigRouteGet(context.Context, string) (ImplResponse, error)
	ConfigRoutesGet(context.Context) (ImplResponse, error)
	ConfigUnmanagedCreate(context.Context, Unmanaged) (ImplResponse, error)
	ConfigUnmanagedDel(context.Context, string) (ImplResponse, error)
	ConfigUnmanagedGet(context.Context, string) (ImplResponse, error)
	ConfigUnmanagedListGet(context.Context) (ImplResponse, error)
}

// SystemApiServicer defines the api actions for the SystemApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type SystemApiServicer interface {
	ConfigGet(context.Context) (ImplResponse, error)
	ConfigPatch(context.Context, Config) (ImplResponse, error)
	ConfigSet(context.Context, Config) (ImplResponse, error)
	PersistConfig(context.Context) (ImplResponse, error)
	ResetConfig(context.Context) (ImplResponse, error)
}
