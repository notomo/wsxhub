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

	if err := filterClause.BatchOperatorType.Validate(); err != nil {
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
	OperatorType      domain.OperatorType `json:"operator"`
	BatchOperatorType domain.OperatorType `json:"batchOperator"`
	Filters           []FilterImpl        `json:"filters"`
	Not               bool                `json:"not"`
}

// Match :
func (clause *FilterClauseImpl) Match(message domain.Message) (bool, error) {
	if len(clause.Filters) == 0 {
		return !clause.Not, nil
	}
	targetMaps := message.Unmarshaled()

	var matched bool
	var err error
	switch clause.OperatorType {
	case domain.OperatorTypeAnd:
		switch clause.BatchOperatorType {
		case domain.OperatorTypeAnd, domain.OperatorTypeDefault:
			matched, err = clause.andMatchAll(targetMaps)
		case domain.OperatorTypeOr:
			matched, err = clause.andMatchOne(targetMaps)
		default:
			return false, errors.New("maybe batch operator type is not validated: " + string(clause.BatchOperatorType))
		}
	case domain.OperatorTypeOr, domain.OperatorTypeDefault:
		switch clause.BatchOperatorType {
		case domain.OperatorTypeAnd, domain.OperatorTypeDefault:
			matched, err = clause.orMatchAll(targetMaps)
		case domain.OperatorTypeOr:
			matched, err = clause.orMatchOne(targetMaps)
		default:
			return false, errors.New("maybe batch operator type is not validated: " + string(clause.BatchOperatorType))
		}
	default:
		return false, errors.New("maybe operator type is not validated: " + string(clause.OperatorType))
	}

	if clause.Not {
		return !matched, err
	}
	return matched, err
}

func (clause *FilterClauseImpl) andMatchAll(targetMaps []map[string]interface{}) (bool, error) {
	for _, target := range targetMaps {
		for _, filter := range clause.Filters {
			matched, err := filter.Match(target)
			if err != nil {
				return false, err
			}
			if !matched {
				return false, nil
			}
		}
	}
	return true, nil
}

func (clause *FilterClauseImpl) andMatchOne(targetMaps []map[string]interface{}) (bool, error) {
	for _, filter := range clause.Filters {
		ok := false
		for _, target := range targetMaps {
			matched, err := filter.Match(target)
			if err != nil {
				return false, err
			}
			if matched {
				ok = true
				break
			}
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (clause *FilterClauseImpl) orMatchAll(targetMaps []map[string]interface{}) (bool, error) {
	for _, target := range targetMaps {
		ok := false
		for _, filter := range clause.Filters {
			matched, err := filter.Match(target)
			if err != nil {
				return false, err
			}
			if matched {
				ok = true
				break
			}
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

func (clause *FilterClauseImpl) orMatchOne(targetMaps []map[string]interface{}) (bool, error) {
	for _, target := range targetMaps {
		for _, filter := range clause.Filters {
			matched, err := filter.Match(target)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
	}
	return false, nil
}

// FilterImpl :
type FilterImpl struct {
	MatchType domain.MatchType       `json:"type"`
	Map       map[string]interface{} `json:"map"`
}

// Match :
func (filter *FilterImpl) Match(targetMap map[string]interface{}) (bool, error) {
	switch filter.MatchType {
	case domain.MatchTypeExact:
		return isSubset(filter.Map, targetMap) && isSubset(targetMap, filter.Map), nil
	case domain.MatchTypeExactKey:
		return isSubsetKey(filter.Map, targetMap) && isSubsetKey(targetMap, filter.Map), nil
	case domain.MatchTypeRegexp:
		return regexpMatch(filter.Map, targetMap), nil
	case domain.MatchTypeContained, domain.MatchTypeDefault:
		return isSubset(filter.Map, targetMap), nil
	case domain.MatchTypeContain:
		return isSubset(targetMap, filter.Map), nil
	case domain.MatchTypeContainedKey:
		return isSubsetKey(filter.Map, targetMap), nil
	case domain.MatchTypeContainKey:
		return isSubsetKey(targetMap, filter.Map), nil
	}
	return false, errors.New("maybe match type is not validated: " + string(filter.MatchType))
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
