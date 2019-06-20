package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/notomo/wsxhub/internal/domain"
)

// FilterClauseFactoryImpl :
type FilterClauseFactoryImpl struct {
}

// FilterClause :
func (factory *FilterClauseFactoryImpl) FilterClause(source string) (domain.FilterClause, error) {
	var filterClause FilterClauseImpl
	if source == "" {
		return &filterClause, nil
	}

	if err := json.Unmarshal([]byte(source), &filterClause); err != nil {
		return nil, err
	}

	if err := filterClause.OperatorType.Validate(); err != nil {
		return nil, err
	}

	for i, filter := range filterClause.Filters {
		filter := filter
		if err := filter.MatchType.Validate(); err != nil {
			return nil, err
		}

		if filter.MatchType != domain.MatchTypeRegexp {
			continue
		}
		regexpMap, err := toRegexpMap(filter.Map)
		if err != nil {
			return nil, err
		}
		filterClause.Filters[i].Map = regexpMap
	}

	return &filterClause, nil
}

func toRegexpMap(filterMap map[string]interface{}) (map[string]interface{}, error) {
	regexpMap := map[string]interface{}{}
	for key, value := range filterMap {
		nestMap, nested := value.(map[string]interface{})
		if nested {
			compiledMap, err := toRegexpMap(nestMap)
			if err != nil {
				return nil, err
			}
			regexpMap[key] = compiledMap
			continue
		}

		pattern, ok := value.(string)
		if !ok {
			msg := fmt.Sprintf("regexp filter values must be string, but actual: %s", value)
			return nil, errors.New(msg)
		}

		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexpMap[key] = compiled
	}

	return regexpMap, nil
}

// FilterClauseImpl :
type FilterClauseImpl struct {
	OperatorType domain.OperatorType `json:"operator"`
	Filters      []FilterImpl        `json:"filters"`
}

// Match :
func (group *FilterClauseImpl) Match(targetMap map[string]interface{}) bool {
	switch group.OperatorType {
	case domain.OperatorTypeAnd:
		return group.andMatch(targetMap)
	case domain.OperatorTypeOr:
		return group.orMatch(targetMap)
	}
	return true
}

func (group *FilterClauseImpl) andMatch(targetMap map[string]interface{}) bool {
	for _, filter := range group.Filters {
		if !filter.Match(targetMap) {
			return false
		}
	}
	return true
}

func (group *FilterClauseImpl) orMatch(targetMap map[string]interface{}) bool {
	for _, filter := range group.Filters {
		if filter.Match(targetMap) {
			return true
		}
	}
	return false
}

// FilterImpl :
type FilterImpl struct {
	MatchType domain.MatchType       `json:"type"`
	Map       map[string]interface{} `json:"map"`
}

// Match :
func (filter *FilterImpl) Match(targetMap map[string]interface{}) bool {
	switch filter.MatchType {
	case domain.MatchTypeExact:
		return isSubset(filter.Map, targetMap) && isSubset(targetMap, filter.Map)
	case domain.MatchTypeExactKey:
		return isSubsetKey(filter.Map, targetMap) && isSubsetKey(targetMap, filter.Map)
	case domain.MatchTypeRegexp:
		return regexpMatch(filter.Map, targetMap)
	case domain.MatchTypeContained:
		return isSubset(filter.Map, targetMap)
	case domain.MatchTypeContain:
		return isSubset(targetMap, filter.Map)
	case domain.MatchTypeContainedKey:
		return isSubsetKey(filter.Map, targetMap)
	case domain.MatchTypeContainKey:
		return isSubsetKey(targetMap, filter.Map)
	}
	return false
}

func regexpMatch(filterMap map[string]interface{}, targetMap map[string]interface{}) bool {
	for key, value := range filterMap {
		targetValue, ok := targetMap[key]
		if !ok {
			return false
		}

		nestMap, nested := value.(map[string]interface{})
		nestTargetMap, nestedTarget := targetValue.(map[string]interface{})
		if nested != nestedTarget {
			return false
		}

		if !nested {
			regexp := value.(*regexp.Regexp)
			targetString, ok := targetValue.(string)
			if !ok {
				return false
			}
			if !regexp.MatchString(targetString) {
				return false
			}
			continue
		}
		if nested && !regexpMatch(nestMap, nestTargetMap) {
			return false
		}
	}
	return true
}

func isSubsetKey(filterMap map[string]interface{}, targetMap map[string]interface{}) bool {
	for key, value := range filterMap {
		targetValue, ok := targetMap[key]
		if !ok {
			return false
		}

		nestMap, nested := value.(map[string]interface{})
		nestTargetMap, nestedTarget := targetValue.(map[string]interface{})
		if nested != nestedTarget {
			return false
		}
		if nested && !isSubsetKey(nestMap, nestTargetMap) {
			return false
		}
	}
	return true
}

func isSubset(filterMap map[string]interface{}, targetMap map[string]interface{}) bool {
	for key, value := range filterMap {
		targetValue, ok := targetMap[key]
		if !ok {
			return false
		}

		nestMap, nested := value.(map[string]interface{})
		nestTargetMap, nestedTarget := targetValue.(map[string]interface{})
		if nested != nestedTarget {
			return false
		}
		if !nested && targetValue != value {
			return false
		}
		if nested && !isSubset(nestMap, nestTargetMap) {
			return false
		}
	}
	return true
}
