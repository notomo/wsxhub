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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := FilterImpl{
				MatchType: domain.MatchTypeExact,
				Map:       test.filter,
			}

			got := f.Match(test.target)

			if got != test.want {
				t.Fatalf("want %v, but %v:", test.want, got)
			}
		})
	}
}
