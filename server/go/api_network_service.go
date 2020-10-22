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
	"net/http"
	"errors"
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

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigGet method not implemented")
}

// ConfigLinkCreate - Configures and brings up a link layer interface 
func (s *NetworkApiService) ConfigLinkCreate(ctx context.Context, link Link) (ImplResponse, error) {
	// TODO - update ConfigLinkCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigLinkCreate method not implemented")
}

// ConfigLinkDel - Brings down and delete a link layer interface 
func (s *NetworkApiService) ConfigLinkDel(ctx context.Context, ifname string) (ImplResponse, error) {
	// TODO - update ConfigLinkDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigLinkDel method not implemented")
}

// ConfigLinkGet - Retrieve link layer interface information 
func (s *NetworkApiService) ConfigLinkGet(ctx context.Context, ifname string) (ImplResponse, error) {
	// TODO - update ConfigLinkGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigLinkGet method not implemented")
}

// ConfigNetNSCreate - Configures an new Network Namespace 
func (s *NetworkApiService) ConfigNetNSCreate(ctx context.Context, netns Netns) (ImplResponse, error) {
	// TODO - update ConfigNetNSCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigNetNSCreate method not implemented")
}

// ConfigNetNSDel - Removes an IP Rule 
func (s *NetworkApiService) ConfigNetNSDel(ctx context.Context, netnsid string) (ImplResponse, error) {
	// TODO - update ConfigNetNSDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigNetNSDel method not implemented")
}

// ConfigNetNSGet - Get a network namespace 
func (s *NetworkApiService) ConfigNetNSGet(ctx context.Context, netnsid string) (ImplResponse, error) {
	// TODO - update ConfigNetNSGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigNetNSGet method not implemented")
}

// ConfigRouteCreate - Configures a route 
func (s *NetworkApiService) ConfigRouteCreate(ctx context.Context, route Route) (ImplResponse, error) {
	// TODO - update ConfigRouteCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRouteCreate method not implemented")
}

// ConfigRouteDel - Brings down and delete an L3 IP route 
func (s *NetworkApiService) ConfigRouteDel(ctx context.Context, routeid string) (ImplResponse, error) {
	// TODO - update ConfigRouteDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRouteDel method not implemented")
}

// ConfigRouteGet - Get a L3 route details 
func (s *NetworkApiService) ConfigRouteGet(ctx context.Context, routeid string) (ImplResponse, error) {
	// TODO - update ConfigRouteGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRouteGet method not implemented")
}

// ConfigRuleCreate - Configures an IP rule 
func (s *NetworkApiService) ConfigRuleCreate(ctx context.Context, body map[string]interface{}) (ImplResponse, error) {
	// TODO - update ConfigRuleCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRuleCreate method not implemented")
}

// ConfigRuleDel - Removes an IP Rule 
func (s *NetworkApiService) ConfigRuleDel(ctx context.Context, ruleid string) (ImplResponse, error) {
	// TODO - update ConfigRuleDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRuleDel method not implemented")
}

// ConfigRuleGet - Get an IP rule details 
func (s *NetworkApiService) ConfigRuleGet(ctx context.Context, ruleid string) (ImplResponse, error) {
	// TODO - update ConfigRuleGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigRuleGet method not implemented")
}

// ConfigSet - Configures and enforces a new live network configuration 
func (s *NetworkApiService) ConfigSet(ctx context.Context, config Config) (ImplResponse, error) {
	// TODO - update ConfigSet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigSet method not implemented")
}

// ConfigVRFCreate - Configures an new VRF 
func (s *NetworkApiService) ConfigVRFCreate(ctx context.Context, body map[string]interface{}) (ImplResponse, error) {
	// TODO - update ConfigVRFCreate with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(201, {}) or use other options such as http.Ok ...
	//return Response(201, nil),nil

	//TODO: Uncomment the next line to return response Response(409, {}) or use other options such as http.Ok ...
	//return Response(409, nil),nil

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigVRFCreate method not implemented")
}

// ConfigVRFDel - Removes a VRF 
func (s *NetworkApiService) ConfigVRFDel(ctx context.Context, vrfid string) (ImplResponse, error) {
	// TODO - update ConfigVRFDel with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigVRFDel method not implemented")
}

// ConfigVRFGet - Get a VRF 
func (s *NetworkApiService) ConfigVRFGet(ctx context.Context, vrfid string) (ImplResponse, error) {
	// TODO - update ConfigVRFGet with the required logic for this service method.
	// Add api_network_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	var err error
	err = nil
	if err != nil {
		return PostErrorResponse(err)
	}
	return Response(http.StatusNotImplemented, nil), errors.New("ConfigVRFGet method not implemented")
}

