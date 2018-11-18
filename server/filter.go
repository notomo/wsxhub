package server

import (
	"encoding/json"
	"fmt"
	"regexp"

	log "github.com/sirupsen/logrus"
)

// StringMapFilter filters values
type StringMapFilter struct {
	stringMap map[string]interface{}
}

// KeyFilter filters keys
type KeyFilter struct {
	stringMap map[string]interface{}
}

// RegexFilter filters values by regular expression
type RegexFilter struct {
	regexMap *RegexMap
}

// NewStringMapFilter creates StringMapFilter
func NewStringMapFilter(stringMap map[string]interface{}) *StringMapFilter {
	return &StringMapFilter{stringMap: stringMap}
}

// NewKeyFilter creates KeyFilter
func NewKeyFilter(stringMap map[string]interface{}) *KeyFilter {
	return &KeyFilter{stringMap: stringMap}
}

// NewRegexFilter creates RegexFilter
func NewRegexFilter(stringMap map[string]interface{}) *RegexFilter {
	return &RegexFilter{regexMap: toRegexMap(stringMap)}
}

// NewStringMapFilterFromString creates StringMapFilter from json string
func NewStringMapFilterFromString(filterString string) (*StringMapFilter, error) {
	stringMap, err := newStringMapFromString(filterString)
	return &StringMapFilter{stringMap: stringMap}, err
}

// NewKeyFilterFromString creates KeyFilter from json string
func NewKeyFilterFromString(filterString string) (*KeyFilter, error) {
	stringMap, err := newStringMapFromString(filterString)
	return &KeyFilter{stringMap: stringMap}, err
}

// NewRegexFilterFromString creates RegexFilter from json string
func NewRegexFilterFromString(filterString string) (*RegexFilter, error) {
	stringMap, err := newStringMapFromString(filterString)
	return &RegexFilter{regexMap: toRegexMap(stringMap)}, err
}

func newStringMapFromString(filterString string) (map[string]interface{}, error) {
	var stringMap interface{}
	if filterString == "" {
		stringMap = map[string]interface{}{}
	} else {
		if err := json.Unmarshal([]byte(filterString), &stringMap); err != nil {
			return nil, err
		}
	}
	value, ok := stringMap.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("filter json parse error")
	}
	return value, nil
}

func toRegexMap(stringMap map[string]interface{}) *RegexMap {
	regexMap := &RegexMap{map[string]RegexMapNode{}}
	for key, value := range stringMap {
		nestMapA, nestedA := value.(map[string]interface{})
		if nestedA {
			regexMap.set(key, toRegexMap(nestMapA))
			continue
		}

		regexString, ok := value.(string)
		if ok {
			regex, err := regexp.Compile(regexString)
			if err != nil {
				log.Warn(err)
				continue
			}
			leaf := &RegexMapLeaf{regex: regex}
			regexMap.set(key, leaf)
			continue
		}
	}
	return regexMap
}

// RegexMapNode represents RegexMap|RegexMapLeaf
type RegexMapNode interface {
	match(e interface{}) bool
}

func (regexMap *RegexMap) match(e interface{}) bool {
	stringMap, isStringMap := e.(map[string]interface{})
	if !isStringMap {
		return false
	}

	for key, node := range regexMap.nodes {
		value, ok := stringMap[key]
		if !ok || !node.match(value) {
			return false
		}
	}

	return true
}

func (regexMap *RegexMap) set(key string, node RegexMapNode) {
	regexMap.nodes[key] = node
}

// RegexMap has RegexMapNodes
type RegexMap struct {
	nodes map[string]RegexMapNode
}

// RegexMapLeaf has compiled regular expression
type RegexMapLeaf struct {
	regex *regexp.Regexp
}

func (regexMapLeaf *RegexMapLeaf) match(e interface{}) bool {
	s, ok := e.(string)
	if !ok {
		return false
	}
	return regexMapLeaf.regex.MatchString(s)
}

// isSubsetOf returns true if stringMap is a subset of this filter
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

// Match returns true if stringMap is filtered
func (filter *KeyFilter) Match(stringMap map[string]interface{}) bool {
	return match(filter.stringMap, stringMap)
}

// Match returns true if values of the stringMap matches filter regular expression
func (filter *RegexFilter) Match(stringMap map[string]interface{}) bool {
	return (*filter.regexMap).match(stringMap)
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
