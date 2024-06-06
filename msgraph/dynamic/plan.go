package dynamic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DynamicSemanticallyEqualFunc func(a, b types.Dynamic) bool

func UseStateWhen(equalFunc DynamicSemanticallyEqualFunc) planmodifier.Dynamic {
	return dynamicUseStateWhen{
		EqualFunc: equalFunc,
	}
}

type dynamicUseStateWhen struct {
	EqualFunc DynamicSemanticallyEqualFunc
}

func (u dynamicUseStateWhen) Description(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u dynamicUseStateWhen) MarkdownDescription(ctx context.Context) string {
	return "Use the state value when new value is functionally equivalent to the old and thus no change is required."
}

func (u dynamicUseStateWhen) PlanModifyDynamic(ctx context.Context, request planmodifier.DynamicRequest, response *planmodifier.DynamicResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	if request.StateValue.IsNull() || request.StateValue.IsUnknown() {
		return
	}
	if semanticallyEqual(ctx, request.ConfigValue, request.StateValue) {
		response.PlanValue = request.StateValue
	}
}

func semanticallyEqual(ctx context.Context, a, b types.Dynamic) bool {
	if a.IsNull() && b.IsNull() {
		return true
	}
	if a.IsNull() || b.IsNull() {
		return false
	}
	aType := a.UnderlyingValue().Type(ctx)
	bType := b.UnderlyingValue().Type(ctx)
	if aType.Equal(types.StringType) && bType.Equal(types.StringType) {
		aJson := a.UnderlyingValue().(types.String).ValueString()
		bJson := b.UnderlyingValue().(types.String).ValueString()
		return normalizeJson(aJson) == normalizeJson(bJson)
	}
	return SemanticallyEqual(a, b)
}

func SemanticallyEqual(a, b types.Dynamic) bool {
	aJson, err := ToJSON(a)
	if err != nil {
		return false
	}
	bJson, err := ToJSON(b)
	if err != nil {
		return false
	}
	return normalizeJson(string(aJson)) == normalizeJson(string(bJson))
}

func normalizeJson(jsonString interface{}) string {
	if jsonString == nil || jsonString == "" {
		return ""
	}
	var j interface{}

	if err := json.Unmarshal([]byte(jsonString.(string)), &j); err != nil {
		return fmt.Sprintf("Error parsing JSON: %+v", err)
	}
	b, _ := json.Marshal(j)
	return string(b)
}
