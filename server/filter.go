package server

import "encoding/json"

// StringMapFilter is
type StringMapFilter struct {
	stringMap map[string]interface{}
}

// NewStringMapFilter is
func NewStringMapFilter(stringMap map[string]interface{}) *StringMapFilter {
	return &StringMapFilter{stringMap: stringMap}
}

// NewStringMapFilterFromString is
func NewStringMapFilterFromString(filterString string) *StringMapFilter {
	var stringMap interface{}
	if filterString == "" {
		stringMap = map[string]interface{}{}
	} else {
		if err := json.Unmarshal([]byte(filterString), &stringMap); err != nil {
			panic(err)
		}
	}
	return &StringMapFilter{stringMap: stringMap.(map[string]interface{})}
}

// isSubsetOf is
func (filter *StringMapFilter) isSubsetOf(stringMap map[string]interface{}) bool {
	for key, value := range filter.stringMap {
		jsonValue, ok := stringMap[key]
		if !ok || jsonValue != value {
			return false
		}
	}
	return true
}
