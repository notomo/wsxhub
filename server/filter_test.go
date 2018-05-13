package server

import "testing"

func TestValueFilter(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		input  S
		filter S
		want   bool
	}{
		{input: S{"otherField": "otherValue"}, filter: S{"neededField": "neededValue"}, want: false},
		{input: S{"neededField": "neededValue"}, filter: S{"neededField": "neededValue"}, want: true},
		{input: S{"neededField": "neededValue", "otherField": "otherValue"}, filter: S{"neededField": "neededValue"}, want: true},
		{input: S{"neededField": "otherValue"}, filter: S{"neededField": "neededValue"}, want: false},
		{input: S{"neededField": "otherValue"}, filter: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, want: false},
		{input: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, filter: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, want: true},
		{input: S{"neededField": S{"nestNeededField": "otherValue"}}, filter: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, want: false},
		{input: S{"neededField": "otherValue"}, filter: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, want: false},
		{input: S{"neededField": S{"nestNeededField": "nestNeededValue", "otherField": "otherValue"}}, filter: S{"neededField": S{"nestNeededField": "nestNeededValue"}}, want: true},
	}

	for _, test := range tests {
		f := NewStringMapFilter(test.filter)
		got := f.isSubsetOf(test.input)
		if got != test.want {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}
