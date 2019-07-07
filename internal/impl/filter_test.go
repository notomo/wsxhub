package impl

import (
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
	"github.com/notomo/wsxhub/internal/mock"
)

func TestMatch(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name         string
		filterClause domain.FilterClause
		targets      []map[string]interface{}
		want         bool
	}{
		{
			name:         "no filter",
			filterClause: &FilterClauseImpl{},
			targets: []S{
				S{"id": "1"},
			},
			want: true,
		},
		{
			name: "no filter with not",
			filterClause: &FilterClauseImpl{
				Not: true,
			},
			targets: []S{
				S{"id": "1"},
			},
			want: false,
		},
		{
			name: "match all filters with all targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeAnd,
				BatchOperatorType: domain.OperatorTypeAnd,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1", "key1": "hoge"},
				S{"key1": "hoge", "id": "1", "key2": "foo"},
			},
			want: true,
		},
		{
			name: "not match all filters with all targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeAnd,
				BatchOperatorType: domain.OperatorTypeAnd,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1", "key1": "hoge"},
				S{"key1": "hoge", "id": "2", "key2": "foo"},
			},
			want: false,
		},
		{
			name: "match all filters with all targets with not",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeAnd,
				BatchOperatorType: domain.OperatorTypeAnd,
				Not:               true,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1", "key1": "hoge"},
				S{"key1": "hoge", "id": "1", "key2": "foo"},
			},
			want: false,
		},
		{
			name: "match all filters with one target",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeAnd,
				BatchOperatorType: domain.OperatorTypeOr,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1", "key1": "hoge", "key2": "foo"},
				S{"id": "2", "key1": "hoge"},
			},
			want: true,
		},
		{
			name: "not match all filters with one target",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeAnd,
				BatchOperatorType: domain.OperatorTypeOr,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "2", "key2": "foo"},
				S{"id": "2"},
			},
			want: false,
		},
		{
			name: "match one filter with all targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeOr,
				BatchOperatorType: domain.OperatorTypeAnd,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "1",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1"},
				S{"key1": "hoge"},
			},
			want: true,
		},
		{
			name: "not match one filter with all targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeOr,
				BatchOperatorType: domain.OperatorTypeAnd,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "2",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "foo",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1"},
				S{"key1": "hoge"},
			},
			want: false,
		},
		{
			name: "match one filter with one targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeOr,
				BatchOperatorType: domain.OperatorTypeOr,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "2",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1"},
				S{"key1": "hoge"},
				S{"id": "3"},
			},
			want: true,
		},
		{
			name: "not match one filter with one targets",
			filterClause: &FilterClauseImpl{
				OperatorType:      domain.OperatorTypeOr,
				BatchOperatorType: domain.OperatorTypeOr,
				Filters: []FilterImpl{
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"id": "2",
						},
					},
					{
						MatchType: domain.MatchTypeContained,
						Map: S{
							"key1": "hoge",
						},
					},
				},
			},
			targets: []S{
				S{"id": "1"},
				S{"key2": "hoge"},
				S{"id": "3"},
			},
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			message := &mock.FakeMessage{
				FakeUnmarshaled: func() []map[string]interface{} {
					return test.targets
				},
			}

			got, err := test.filterClause.Match(message)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}

	t.Run("invalid operator", func(t *testing.T) {
		message := &mock.FakeMessage{
			FakeUnmarshaled: func() []map[string]interface{} {
				return []S{}
			},
		}

		filterClause := &FilterClauseImpl{
			OperatorType: domain.OperatorType("invalid"),
			Filters: []FilterImpl{
				{
					Map: S{},
				},
			},
		}

		if _, err := filterClause.Match(message); err == nil {
			t.Fatalf("should be error")
		}
	})

}

func TestExactMatch(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "no nest, other key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, the same key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
		{
			name:   "no nest, other value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "otherValue"},
			want:   false,
		},
		{
			name:   "nest filter, no nest target",
			filter: S{"neededKey": S{"nestNeededKey": "nestNeededValue"}},
			target: S{"neededKey": "otherValue"},
			want:   false,
		},
		{
			name:   "nest filter, nest target, the same key and value",
			filter: S{"neededKey": S{"nestNeededKey": "nestNeededValue"}},
			target: S{"neededKey": S{"nestNeededKey": "nestNeededValue"}},
			want:   true,
		},
		{
			name:   "nest filter, nest target, other value",
			filter: S{"neededKey": S{"nestNeededKey": "nestNeededValue"}},
			target: S{"neededKey": S{"nestNeededKey": "otherValue"}},
			want:   false,
		},
		{
			name:   "no nest filter, nest target",
			filter: S{"neededKey": "otherValue"},
			target: S{"neededKey": S{"nestNeededKey": "nestNeededValue"}},
			want:   false,
		},
		{
			name:   "no nest filter, slice target",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": []string{"neededValue"}},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeExact,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestExactKeyMatch(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "no nest, other key",
			filter: S{"neededKey": "value"},
			target: S{"otherKey": "value"},
			want:   false,
		},
		{
			name:   "no nest, the same key",
			filter: S{"neededKey": "value"},
			target: S{"neededKey": "otherValue"},
			want:   true,
		},
		{
			name:   "nest filter, no nest target",
			filter: S{"neededKey": S{"nestNeededKey": "value"}},
			target: S{"neededKey": "value"},
			want:   false,
		},
		{
			name:   "nest filter, nest target, other key",
			filter: S{"neededKey": S{"nestNeededKey": "value"}},
			target: S{"neededKey": S{"nestOtherKey": "otherValue"}},
			want:   false,
		},
		{
			name:   "nest filter, nest target, the same key",
			filter: S{"neededKey": S{"nestNeededKey": "value"}},
			target: S{"neededKey": S{"nestNeededKey": "otherValue"}},
			want:   true,
		},
		{
			name:   "no nest filter, nest target",
			filter: S{"neededKey": "value"},
			target: S{"neededKey": S{"nestNeededKey": "value"}},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeExactKey,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestRegexpMatch(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "no nest, other key",
			filter: S{"neededKey": "value"},
			target: S{"otherKey": "value"},
			want:   false,
		},
		{
			name:   "no nest, not matched value",
			filter: S{"neededKey": "value"},
			target: S{"neededKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, matched value",
			filter: S{"neededKey": ".*"},
			target: S{"neededKey": "value"},
			want:   true,
		},
		{
			name:   "no nest, not string target",
			filter: S{"neededKey": ".*"},
			target: S{"neededKey": 8888},
			want:   false,
		},
		{
			name:   "nest filter, no nest target",
			filter: S{"neededKey": S{"nestNeededKey": "value"}},
			target: S{"neededKey": "value"},
			want:   false,
		},
		{
			name:   "nest filter, nest target, not matched value",
			filter: S{"neededKey": S{"nestNeededKey": "value"}},
			target: S{"neededKey": S{"nestNeededKey": "otherValue"}},
			want:   false,
		},
		{
			name:   "nest filter, nest target, matched value",
			filter: S{"neededKey": S{"nestNeededKey": "otherValue|value"}},
			target: S{"neededKey": S{"nestNeededKey": "value"}},
			want:   true,
		},
		{
			name:   "no nest filter, nest target",
			filter: S{"neededKey": "value"},
			target: S{"neededKey": S{"nestNeededKey": "value"}},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			regexpMap, err := toRegexpMap(test.filter)
			if err != nil {
				t.Fatalf(err.Error())
			}
			f := FilterImpl{
				MatchType: domain.MatchTypeRegexp,
				Map:       regexpMap,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestContained(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "empty filter, target",
			filter: S{},
			target: S{},
			want:   true,
		},
		{
			name:   "empty filter",
			filter: S{},
			target: S{"otherKey": "value"},
			want:   true,
		},
		{
			name:   "empty target",
			filter: S{"neededKey": "neededValue"},
			target: S{},
			want:   false,
		},
		{
			name:   "no nest, other key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, the same key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
		{
			name:   "no nest, contained key value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			want:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeContained,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestContain(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "empty filter, target",
			filter: S{},
			target: S{},
			want:   true,
		},
		{
			name:   "empty filter",
			filter: S{},
			target: S{"otherKey": "value"},
			want:   false,
		},
		{
			name:   "empty target",
			filter: S{"neededKey": "neededValue"},
			target: S{},
			want:   true,
		},
		{
			name:   "no nest, other key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, the same key and value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
		{
			name:   "no nest, contained key value",
			filter: S{"neededKey": "neededValue"},
			target: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, contain key value",
			filter: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeContain,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestContainedKey(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "no nest, other key",
			filter: S{"neededKey": "neededValue"},
			target: S{"otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, the same key",
			filter: S{"neededKey": "otherValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
		{
			name:   "no nest, contained key",
			filter: S{"neededKey": "otherValue"},
			target: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			want:   true,
		},
		{
			name:   "no nest, contain key",
			filter: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			target: S{"neededKey": "otherValue"},
			want:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeContainedKey,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}

func TestContainKey(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		name   string
		target S
		filter S
		want   bool
	}{
		{
			name:   "no nest, other key",
			filter: S{"neededKey": "neededValue"},
			target: S{"otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, the same key",
			filter: S{"neededKey": "otherValue"},
			target: S{"neededKey": "neededValue"},
			want:   true,
		},
		{
			name:   "no nest, contained key",
			filter: S{"neededKey": "otherValue"},
			target: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			want:   false,
		},
		{
			name:   "no nest, contain key",
			filter: S{"neededKey": "neededValue", "otherKey": "otherValue"},
			target: S{"neededKey": "otherValue"},
			want:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeContainKey,
				Map:       test.filter,
			}

			got, err := f.Match(test.target)
			if err != nil {
				t.Fatalf("should not be error: %v", err)
			}

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}

	t.Run("invalid match type", func(t *testing.T) {
		f := FilterImpl{
			MatchType: domain.MatchType("invalid"),
			Map:       S{},
		}

		if _, err := f.Match(S{}); err == nil {
			t.Fatalf("should be error")
		}
	})
}
