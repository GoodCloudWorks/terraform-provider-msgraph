package id

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ID struct {
	value      *url.URL
	apiVersion string
	Path       string
}

func New(collection string, objectId string) *ID {
	path := strings.Trim(fmt.Sprintf("%s/%s", collection, objectId), "/")
	return &ID{
		value: &url.URL{
			Path: path,
		},
		Path: path,
	}
}

func (id *ID) AsString() types.String {
	return types.StringValue(id.Path)
}

func Parse(value string) (*ID, error) {
	url, err := url.Parse(strings.Trim(value, "/"))
	if err != nil {
		return nil, err
	}

	apiVersion := ""
	if strings.HasPrefix(url.Path, "v1.0/") || strings.HasPrefix(url.Path, "beta/") {
		apiVersion = url.Path[:4]
		url, err = url.Parse(url.Path[5:])
		if err != nil {
			return nil, err
		}
	}

	if url.Path == "" {
		return nil, fmt.Errorf("invalid id: %s", value)
	}

	return &ID{
		value:      url,
		apiVersion: apiVersion,
		Path:       url.Path,
	}, nil
}

func ParseString(value types.String) (*ID, error) {
	return Parse(value.ValueString())
}

func (id *ID) Collection() string {
	index := strings.LastIndex(id.value.Path, "/")

	if index == -1 {
		return ""
	}

	return id.value.Path[:index]
}

func (id *ID) ObjectId() string {
	index := strings.LastIndex(id.value.Path, "/")
	return id.value.Path[index+1:]
}

func (id *ID) ApiVersion() string {
	return id.apiVersion
}
