package dynamic

import (
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UpdateWithSchemaPreservation(source, target types.Dynamic) (types.Dynamic, error) {
	if source.IsNull() || target.IsNull() {
		return target, nil
	}

	sourceJSON, err := ToJSON(source)
	if err != nil {
		return types.DynamicNull(), err
	}

	targetJSON, err := ToJSON(target)
	if err != nil {
		return types.DynamicNull(), err
	}

	var sourceObject interface{}
	if err := json.Unmarshal(sourceJSON, &sourceObject); err != nil {
		return types.DynamicNull(), err
	}

	var targetObject interface{}
	if err := json.Unmarshal(targetJSON, &targetObject); err != nil {
		return types.DynamicNull(), err
	}

	resultObject := updateObjects(sourceObject, targetObject)

	resultJSON, err := json.Marshal(resultObject)
	if err != nil {
		return types.DynamicNull(), err
	}

	return FromJSONImplied(resultJSON)
}

func updateObjects(source, target interface{}) interface{} {
	switch target := target.(type) {
	case map[string]interface{}:
		source, ok := source.(map[string]interface{})
		if !ok {
			return source
		}

		return updateMap(source, target)

	case []interface{}:
		source, ok := source.([]interface{})
		if !ok {
			return target
		}

		return updateArray(source, target)
	}

	return target
}

func updateMap(source, target map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key := range target {
		if sourceValue, ok := source[key]; ok {
			result[key] = updateObjects(sourceValue, target[key])
		}
	}

	return result
}

func updateArray(source, target []interface{}) []interface{} {
	result := make([]interface{}, 0, len(source))

	for i := range target {
		if i < len(source) {
			result = append(result, updateObjects(source[i], target[i]))
		}
	}

	for i := len(target); i < len(source); i++ {
		result = append(result, source[i])
	}

	return result
}
