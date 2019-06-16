package impl

import (
	"encoding/json"

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

	return &filterClause, nil
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
		return filter.exactMatch(targetMap)
	case domain.MatchTypeExactKey:
		return filter.exactKeyMatch(targetMap)
	case domain.MatchTypeRegexp:
		return filter.regexpMatch(targetMap)
	}
	return false
}

func (filter *FilterImpl) exactMatch(targetMap map[string]interface{}) bool {
	return isSubset(filter.Map, targetMap) && isSubset(targetMap, filter.Map)
}

func (filter *FilterImpl) exactKeyMatch(targetMap map[string]interface{}) bool {
	return isSubsetKey(filter.Map, targetMap) && isSubsetKey(targetMap, filter.Map)
}

func (filter *FilterImpl) regexpMatch(targetMap map[string]interface{}) bool {
	return false
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
		if nested && !isSubset(nestMap, nestTargetMap) {
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
