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

func TestKeyFilter(t *testing.T) {
	type S = map[string]interface{}

	tests := []struct {
		input  S
		filter S
		want   bool
	}{
		{input: S{"otherField": "value"}, filter: S{"field": true}, want: false},
		{input: S{"field": "value"}, filter: S{"field": true}, want: true},
		{input: S{"field": "value"}, filter: S{"field": false}, want: false},
		{input: S{}, filter: S{"field": false}, want: true},
		{input: S{"field": "value", "otherField": "value"}, filter: S{"field": true}, want: true},
		{input: S{"field": S{"nestField": "value"}}, filter: S{"field": S{"nestField": true}}, want: true},
		{input: S{"field": S{"nestField": "value"}}, filter: S{"field": S{"nestField": false}}, want: false},
		{input: S{"field": S{"nestField": "value"}}, filter: S{"field": S{"nestOtherField": false}}, want: true},
		{input: S{"field": "value"}, filter: S{"field": S{"nestField": true}}, want: false},
		{input: S{"field": "value"}, filter: S{"field": S{"nestField": false}}, want: false},
		{input: S{}, filter: S{"field": S{"nestField": false}}, want: true},
		{input: S{}, filter: S{"field": S{"nestField": false, "nestField2": true}}, want: false},
		{input: S{"field": S{"nestField2": S{"nest2Field": "value"}}}, filter: S{"field": S{"nestField": false, "nestField2": S{"nest2Field": true}}}, want: true},
		{input: S{}, filter: S{"field": S{"nestField": false, "nestField2": S{"nest2Field": true}}}, want: false},
		{input: S{}, filter: S{"field": S{"nestField": false, "nestField2": S{"nest2Field": false}}}, want: true},
	}

	for _, test := range tests {
		f := NewKeyFilter(test.filter)
		got := f.Match(test.input)
		if got != test.want {
			t.Fatalf("input %q, filter %q, want %q, but %q:", test.input, test.filter, test.want, got)
		}
	}
}
