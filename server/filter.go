package server

// StringMapFilter is
type StringMapFilter struct {
	stringMap map[string]interface{}
}

// NewStringMapFilter is
func NewStringMapFilter(stringMap map[string]interface{}) *StringMapFilter {
	return &StringMapFilter{stringMap: stringMap}
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
