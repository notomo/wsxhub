package mock

import (
	"github.com/notomo/wsxhub/internal/domain"
)

// FakeFilterClause :
type FakeFilterClause struct {
	domain.FilterClause
	FakeMatch func(map[string]interface{}) bool
}

// Match :
func (filterClause *FakeFilterClause) Match(targetMap map[string]interface{}) bool {
	return filterClause.FakeMatch(targetMap)
}
