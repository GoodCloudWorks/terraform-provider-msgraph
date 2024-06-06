package id

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ID struct {
	value *url.URL
	Path  string
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

	return &ID{
		value: url,
		Path:  url.Path,
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
	return id.value.Query().Get("api-version")
}
