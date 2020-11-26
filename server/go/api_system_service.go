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
)

// SystemApiService is a service that implents the logic for the SystemApiServicer
// This service should implement the business logic for every endpoint for the SystemApi API.
// Include any external packages or services that will be required by this service.
type SystemApiService struct {
}

// NewSystemApiService creates a default api service
func NewSystemApiService() SystemApiServicer {
	return &SystemApiService{}
}

// ConfigGet - Get current live configuration 
func (s *SystemApiService) ConfigGet(ctx context.Context) (ImplResponse, error) {
	// TODO - update ConfigGet with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ConfigGet method not implemented")
	return GetErrorResponse(err, nil)
}

// ConfigPatch - Patch existing configuration with new one 
func (s *SystemApiService) ConfigPatch(ctx context.Context, config Config) (ImplResponse, error) {
	// TODO - update ConfigPatch with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ConfigPatch method not implemented")
	return PatchErrorResponse(err, nil)
}

// ConfigSet - Replace existing configuration with new one 
func (s *SystemApiService) ConfigSet(ctx context.Context, config Config) (ImplResponse, error) {
	// TODO - update ConfigSet with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	//return Response(400, nil),nil

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ConfigSet method not implemented")
	return PutErrorResponse(err, nil)
}

// PersistConfig - Persist live configuration
func (s *SystemApiService) PersistConfig(ctx context.Context) (ImplResponse, error) {
	// TODO - update PersistConfig with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("PersistConfig method not implemented")
	return PostErrorResponse(err, nil)
}

// ResetConfig - Reload persisted configuration back
func (s *SystemApiService) ResetConfig(ctx context.Context) (ImplResponse, error) {
	// TODO - update ResetConfig with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("ResetConfig method not implemented")
	return PostErrorResponse(err, nil)
}

// RunDiagnostics - Run a diagnostic session
func (s *SystemApiService) RunDiagnostics(ctx context.Context) (ImplResponse, error) {
	// TODO - update RunDiagnostics with the required logic for this service method.
	// Add api_system_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	err := errors.New("RunDiagnostics method not implemented")
	return PostErrorResponse(err, nil)
}

