/*
 * netConfD API
 *
 * Network Configurator service
 *
 * API version: 0.1.0
 * Contact: support@athonet.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"context"
	"errors"

	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/logger"
	"gitlab.lan.athonet.com/riccardo.manfrin/netconfd/nc"
)

// NetworkApiService is a service that implents the logic for the NetworkApiServicer
// This service should implement the business logic for every endpoint for the NetworkApi API.
// Include any external packages or services that will be required by this service.
type NetworkApiService struct {
}

// NewNetworkApiService creates a default api service
func NewNetworkApiService() NetworkApiServicer {
	return &NetworkApiService{}
}

// ConfigGet - Configures and enforces a new live network configuration
func (s *NetworkApiService) ConfigGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ConfigGet method not implemented")
	return GetErrorResponse(err, nil)
}

func ncLinkFormat(link Link) (nc.Link, error) {

	nclink := nc.Link{
		Ifname:   link.GetIfname(),
		Linkinfo: nc.LinkLinkinfo{InfoKind: *link.GetLinkinfo().InfoKind},
		Mtu:      link.GetMtu(),
		LinkType: link.GetLinkType(),
	}
	switch nclink.Linkinfo.InfoKind {
	case "bond":
		{
			nclink.Linkinfo.InfoData.Mode = *link.GetLinkinfo().InfoData.Mode
		}
	}
	return nclink, nil
}

// ConfigLinkCreate - Configures and brings up a link layer interface
func (s *NetworkApiService) ConfigLinkCreate(ctx context.Context, link Link) (ImplResponse, error) {
	// TODO - update ConfigLinkCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []string{}) or use other options such as http.Ok ...
	//return Response(201, []string{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	nclink, err := ncLinkFormat(link)
	if err != nil {
		return PostErrorResponse(err, nil)
	}
	err = nc.LinkCreate(nclink)
	return PostErrorResponse(err, nil)
}

// ConfigLinkDel - Brings down and delete a link layer interface
func (s *NetworkApiService) ConfigLinkDel(ctx context.Context, ifname string) (ImplResponse, error) {
	// TODO - update ConfigLinkDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigLinkDel method not implemented")
	return DeleteErrorResponse(err, nil)
}
func ncLinkParse(nclink nc.Link) Link {
	link := Link{
		Id:       &nclink.Ifname,
		Ifname:   nclink.Ifname,
		Ifindex:  &nclink.Ifindex,
		LinkType: nclink.LinkType,
	}
	if len(nclink.AddrInfo) > 0 {
		lai := make([]LinkAddrInfo, len(nclink.AddrInfo))
		for i, a := range nclink.AddrInfo {
			ip := a.Local.Address()
			lai[i].Local = &Ip{string: &ip}
		}
		link.AddrInfo = &lai
	}

	lli := LinkLinkinfo{}
	if nclink.Linkinfo.InfoKind != "" {
		lli.InfoKind = &nclink.Linkinfo.InfoKind
		link.Linkinfo = &lli
	}
	isd := LinkLinkinfoInfoSlaveData{}
	lli.InfoSlaveData = &isd

	switch nclink.Linkinfo.InfoKind {
	case "bond":
		{
			id := LinkLinkinfoInfoData{}
			id.Mode = &nclink.Linkinfo.InfoData.Mode
			id.Miimon = &nclink.Linkinfo.InfoData.Miimon
			id.Updelay = &nclink.Linkinfo.InfoData.Updelay
			id.Downdelay = &nclink.Linkinfo.InfoData.Downdelay
			lli.InfoData = &id
		}
	case "device":
	case "bridge":
	case "dummy":
	case "ppp":
	default:
		{
			logger.Log.Warning("Unknown Link Kind")
		}
	}
	if nclink.Master != "" {
		icisd := &nclink.Linkinfo.InfoSlaveData
		link.SetMaster(nclink.Master)
		isd.SetState(icisd.State)
		isd.SetLinkFailureCount(int32(icisd.LinkFailureCount))
		isd.SetMiiStatus(icisd.MiiStatus)
		isd.SetPermHwaddr(icisd.PermHwaddr)
	}

	return link
}

// ConfigLinkGet - Retrieve link layer interface information
func (s *NetworkApiService) ConfigLinkGet(ctx context.Context, ifname string) (ImplResponse, error) {
	// TODO - update ConfigLinkGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	nclink, err := nc.LinkGet(nc.LinkID(ifname))
	if err != nil {
		return GetErrorResponse(err, nil)
	}

	return GetErrorResponse(err, ncLinkParse(nclink))
}

// ConfigLinksGet - Get all link layer interfaces
func (s *NetworkApiService) ConfigLinksGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigLinksGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Link{}) or use other options such as http.Ok ...
	//return Response(200, []Link{}), nil
	var links []Link
	nclinks, err := nc.LinksGet()
	if err == nil {
		links = make([]Link, len(nclinks))
		for i, l := range nclinks {
			links[i] = ncLinkParse(l)
		}
	}

	return GetErrorResponse(err, links)
}

// ConfigNFTableCreate - Configures an new NFTable
func (s *NetworkApiService) ConfigNFTableCreate(ctx context.Context, body map[string]interface{}) (ImplResponse, error) {
	// TODO - update ConfigNFTableCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []int32{}) or use other options such as http.Ok ...
	//return Response(201, []int32{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	err := errors.New("ConfigNFTableCreate method not implemented")
	return PostErrorResponse(err, nil)
}

// ConfigNFTableDel - Removes a NFTable
func (s *NetworkApiService) ConfigNFTableDel(ctx context.Context, nftableid int32) (ImplResponse, error) {
	// TODO - update ConfigNFTableDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigNFTableDel method not implemented")
	return DeleteErrorResponse(err, nil)
}

// ConfigNFTableGet - Get a NFTable
func (s *NetworkApiService) ConfigNFTableGet(ctx context.Context, nftableid int32) (ImplResponse, error) {
	// TODO - update ConfigNFTableGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigNFTableGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigNFTablesGet - Get the list all NFTables
func (s *NetworkApiService) ConfigNFTablesGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigNFTablesGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []int32{}) or use other options such as http.Ok ...
	//return Response(200, []int32{}), nil

	err := errors.New("ConfigNFTablesGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigNetNSCreate - Configures an new Network Namespace
func (s *NetworkApiService) ConfigNetNSCreate(ctx context.Context, netns Netns) (ImplResponse, error) {
	// TODO - update ConfigNetNSCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []string{}) or use other options such as http.Ok ...
	//return Response(201, []string{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	err := errors.New("ConfigNetNSCreate method not implemented")
	return PostErrorResponse(err, nil)
}

// ConfigNetNSDel - Removes an IP Rule
func (s *NetworkApiService) ConfigNetNSDel(ctx context.Context, netnsid string) (ImplResponse, error) {
	// TODO - update ConfigNetNSDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigNetNSDel method not implemented")
	return DeleteErrorResponse(err, nil)
}

// ConfigNetNSGet - Get a network namespace
func (s *NetworkApiService) ConfigNetNSGet(ctx context.Context, netnsid string) (ImplResponse, error) {
	// TODO - update ConfigNetNSGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigNetNSGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigNetNSsGet - Get the list all network namespaces
func (s *NetworkApiService) ConfigNetNSsGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigNetNSsGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []string{}) or use other options such as http.Ok ...
	//return Response(200, []string{}), nil

	err := errors.New("ConfigNetNSsGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigRouteCreate - Configures a route
func (s *NetworkApiService) ConfigRouteCreate(ctx context.Context, route Route) (ImplResponse, error) {
	// TODO - update ConfigRouteCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []int32{}) or use other options such as http.Ok ...
	//return Response(201, []int32{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	err := errors.New("ConfigRouteCreate method not implemented")
	return PostErrorResponse(err, nil)
}

// ConfigRouteDel - Brings down and delete an L3 IP route
func (s *NetworkApiService) ConfigRouteDel(ctx context.Context, routeid int32) (ImplResponse, error) {
	// TODO - update ConfigRouteDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigRouteDel method not implemented")
	return DeleteErrorResponse(err, nil)
}

// ConfigRouteGet - Get a L3 route details
func (s *NetworkApiService) ConfigRouteGet(ctx context.Context, routeid int32) (ImplResponse, error) {
	// TODO - update ConfigRouteGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigRouteGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigRoutesGet - Get all routing table routes
func (s *NetworkApiService) ConfigRoutesGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigRoutesGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []Route{}) or use other options such as http.Ok ...
	//return Response(200, []Route{}), nil

	routes, err := nc.RoutesGet()
	return GetErrorResponse(err, routes)
}

// ConfigRuleCreate - Configures an IP rule
func (s *NetworkApiService) ConfigRuleCreate(ctx context.Context, body map[string]interface{}) (ImplResponse, error) {
	// TODO - update ConfigRuleCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []int32{}) or use other options such as http.Ok ...
	//return Response(201, []int32{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	err := errors.New("ConfigRuleCreate method not implemented")
	return PostErrorResponse(err, nil)
}

// ConfigRuleDel - Removes an IP Rule
func (s *NetworkApiService) ConfigRuleDel(ctx context.Context, ruleid int32) (ImplResponse, error) {
	// TODO - update ConfigRuleDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigRuleDel method not implemented")
	return DeleteErrorResponse(err, nil)
}

// ConfigRuleGet - Get an IP rule details
func (s *NetworkApiService) ConfigRuleGet(ctx context.Context, ruleid int32) (ImplResponse, error) {
	// TODO - update ConfigRuleGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigRuleGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigRulesGet - Get all ip rules list
func (s *NetworkApiService) ConfigRulesGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigRulesGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []int32{}) or use other options such as http.Ok ...
	//return Response(200, []int32{}), nil

	err := errors.New("ConfigRulesGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigSet - Configures and enforces a new live network configuration
func (s *NetworkApiService) ConfigSet(ctx context.Context, config Config) (ImplResponse, error) {
	// TODO - update ConfigSet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ConfigSet method not implemented")
	return PutErrorResponse(err, nil)
}

// ConfigVRFCreate - Configures an new VRF
func (s *NetworkApiService) ConfigVRFCreate(ctx context.Context, body map[string]interface{}) (ImplResponse, error) {
	// TODO - update ConfigVRFCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, []int32{}) or use other options such as http.Ok ...
	//return Response(201, []int32{}), nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	err := errors.New("ConfigVRFCreate method not implemented")
	return PostErrorResponse(err, nil)
}

// ConfigVRFDel - Removes a VRF
func (s *NetworkApiService) ConfigVRFDel(ctx context.Context, vrfid int32) (ImplResponse, error) {
	// TODO - update ConfigVRFDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigVRFDel method not implemented")
	return DeleteErrorResponse(err, nil)
}

// ConfigVRFGet - Get a VRF
func (s *NetworkApiService) ConfigVRFGet(ctx context.Context, vrfid int32) (ImplResponse, error) {
	// TODO - update ConfigVRFGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	err := errors.New("ConfigVRFGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigVRFsGet - Get the list all VRFs
func (s *NetworkApiService) ConfigVRFsGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigVRFsGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, []int32{}) or use other options such as http.Ok ...
	//return Response(200, []int32{}), nil

	err := errors.New("ConfigVRFsGet method not implemented")
	return GetErrorResponse(err, nil)
}
