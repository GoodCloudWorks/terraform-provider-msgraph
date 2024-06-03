package client

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MsGraphClient interface {
	GetToken(context context.Context) (string, error)
	R(context context.Context, apiVersion types.String) *resty.Request
	URL(id types.String) string
}
