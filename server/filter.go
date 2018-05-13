package server

import (
	"encoding/json"
)

// StringMapFilter is
type StringMapFilter struct {
	stringMap map[string]interface{}
}

// KeyFilter is
type KeyFilter struct {
	stringMap map[string]interface{}
}

// NewStringMapFilter is
func NewStringMapFilter(stringMap map[string]interface{}) *StringMapFilter {
	return &StringMapFilter{stringMap: stringMap}
}

// NewKeyFilter is
func NewKeyFilter(stringMap map[string]interface{}) *KeyFilter {
	return &KeyFilter{stringMap: stringMap}
}

// NewStringMapFilterFromString is
func NewStringMapFilterFromString(filterString string) *StringMapFilter {
	return &StringMapFilter{stringMap: newStringMapFromString(filterString)}
}

// NewKeyFilterFromString is
func NewKeyFilterFromString(filterString string) *KeyFilter {
	return &KeyFilter{stringMap: newStringMapFromString(filterString)}
}

func newStringMapFromString(filterString string) map[string]interface{} {
	var stringMap interface{}
	if filterString == "" {
		stringMap = map[string]interface{}{}
	} else {
		if err := json.Unmarshal([]byte(filterString), &stringMap); err != nil {
			panic(err)
		}
	}
	value, ok := stringMap.(map[string]interface{})
	if !ok {
		panic(stringMap)
	}
	return value
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
		if nestedA && !isSubset(nestMapA, nestMapB) {
			return false
		}
	}
	return true
}

// Match is
func (filter *KeyFilter) Match(stringMap map[string]interface{}) bool {
	return match(filter.stringMap, stringMap)
}

func match(a map[string]interface{}, b map[string]interface{}) bool {
	for key, value := range a {
		jsonValue, ok := b[key]
		nestMapA, nestedA := value.(map[string]interface{})
		if !ok && nestedA {
			if !isAllFalse(nestMapA) {
				return false
			}
			continue
		}
		nestMapB, nestedB := jsonValue.(map[string]interface{})
		if nestedA != nestedB {
			return false
		}
		if !nestedA && ok != value.(bool) {
			return false
		}
		if nestedA && !match(nestMapA, nestMapB) {
			return false
		}
	}
	return true
}

func isAllFalse(m map[string]interface{}) bool {
	for _, value := range m {
		if nestMap, nested := value.(map[string]interface{}); nested {
			if !isAllFalse(nestMap) {
				return false
			}
			continue
		}
		if value.(bool) {
			return false
		}
	}
	return true
}
