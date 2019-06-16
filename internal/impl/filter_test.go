package impl

import (
	"testing"

	"github.com/notomo/wsxhub/internal/domain"
)

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

			got := f.Match(test.target)

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

			got := f.Match(test.target)

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

			got := f.Match(test.target)

			if got != test.want {
				t.Errorf("want %v, but %v:", test.want, got)
			}
		})
	}
}
