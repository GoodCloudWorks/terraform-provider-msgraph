package client

import (
	"context"
)

const (
	METHOD_GET    = "GET"
	METHOD_POST   = "POST"
	METHOD_PATCH  = "PATCH"
	METHOD_PUT    = "PUT"
	METHOD_DELETE = "DELETE"
)

type MsGraphRequest struct {
	client     *MsGraphClient
	method     string
	path       string
	apiVersion string
}

func (client *MsGraphClient) Request(path string) *MsGraphRequest {
	return &MsGraphRequest{
		client:     client,
		path:       path,
		method:     METHOD_GET,
		apiVersion: client.Options.ApiVersion,
	}
}

func (request *MsGraphRequest) Method(method string) *MsGraphRequest {
	request.method = method
	return request
}

func (request *MsGraphRequest) ApiVersion(apiVersion string) *MsGraphRequest {
	request.apiVersion = apiVersion
	return request
}

func (request *MsGraphRequest) Get(ctx context.Context) (*interface{}, error) {
	path := request.apiVersion + "/" + request.path
	return request.client.Get(ctx, path)
}
