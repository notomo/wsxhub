package mock

import (
	"github.com/notomo/wsxhub/internal/domain"
)

// FakeFilterClause :
type FakeFilterClause struct {
	domain.FilterClause
	FakeMatch func(domain.Message) bool
}

// Match :
func (filterClause *FakeFilterClause) Match(message domain.Message) bool {
	return filterClause.FakeMatch(message)
}
