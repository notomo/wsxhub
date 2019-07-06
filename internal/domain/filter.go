package domain

import "errors"

// OperatorType :
type OperatorType string

var (
	// OperatorTypeAnd :
	OperatorTypeAnd = OperatorType("and")
	// OperatorTypeOr :
	OperatorTypeOr = OperatorType("or")
	// OperatorTypeDefault :
	OperatorTypeDefault = OperatorType("")
)

// Validate :
func (operatorType OperatorType) Validate() error {
	value := string(operatorType)
	for _, typ := range operatorTypes() {
		if value == string(typ) {
			return nil
		}
	}
	return errors.New("invalid OperatorType: " + value)
}

func operatorTypes() []OperatorType {
	return []OperatorType{
		OperatorTypeAnd,
		OperatorTypeOr,
		OperatorTypeDefault,
	}
}

// MatchType :
type MatchType string

var (
	// MatchTypeExact :
	MatchTypeExact = MatchType("exact")
	// MatchTypeExactKey :
	MatchTypeExactKey = MatchType("exactKey")
	// MatchTypeRegexp :
	MatchTypeRegexp = MatchType("regexp")
	// MatchTypeContained :
	MatchTypeContained = MatchType("contained")
	// MatchTypeContain :
	MatchTypeContain = MatchType("contain")
	// MatchTypeContainedKey :
	MatchTypeContainedKey = MatchType("containedKey")
	// MatchTypeContainKey :
	MatchTypeContainKey = MatchType("containKey")
	// MatchTypeDefault :
	MatchTypeDefault = MatchType("")
)

// Validate :
func (matchType MatchType) Validate() error {
	value := string(matchType)
	for _, typ := range matchTypes() {
		if value == string(typ) {
			return nil
		}
	}
	return errors.New("invalid MatchType: " + value)
}

func matchTypes() []MatchType {
	return []MatchType{
		MatchTypeExact,
		MatchTypeExactKey,
		MatchTypeRegexp,
		MatchTypeContained,
		MatchTypeContain,
		MatchTypeContainedKey,
		MatchTypeContainKey,
		MatchTypeDefault,
	}
}

// FilterClauseFactory :
type FilterClauseFactory interface {
	FilterClause(string) (FilterClause, error)
}

// FilterClause :
type FilterClause interface {
	Match(Message) (bool, error)
}
