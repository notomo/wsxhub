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
	return isSubset(filter.stringMap, stringMap)
}

func isSubset(a map[string]interface{}, b map[string]interface{}) bool {
	for key, value := range a {
		jsonValue, ok := b[key]
		if !ok {
			return false
		}
		nestMapA, nestedA := value.(map[string]interface{})
		nestMapB, nestedB := jsonValue.(map[string]interface{})
		if nestedA != nestedB {
			return false
		}
		if !nestedA && jsonValue != value {
			return false
		}
		if nestedA {
			return isSubset(nestMapA, nestMapB)
		}
	}
	return true
}
