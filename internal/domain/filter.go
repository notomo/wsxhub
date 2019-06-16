package domain

// OperatorType :
type OperatorType string

var (
	// OperatorTypeAnd :
	OperatorTypeAnd = OperatorType("and")
	// OperatorTypeOr :
	OperatorTypeOr = OperatorType("or")
)

// MatchType :
type MatchType string

var (
	// MatchTypeExact :
	MatchTypeExact = MatchType("exact")
	// MatchTypeExactKey :
	MatchTypeExactKey = MatchType("exactKey")
	// MatchTypeRegexp :
	MatchTypeRegexp = MatchType("regexp")
)

// FilterClauseFactory :
type FilterClauseFactory interface {
	FilterClause(string) (FilterClause, error)
}

// FilterClause :
type FilterClause interface {
	Match(map[string]interface{}) bool
}
