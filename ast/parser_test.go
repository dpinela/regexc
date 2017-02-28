package ast

import (
	"reflect"
	"testing"
)

type testCase struct {
	Name  string
	Regex string
	Want  Node
}

var testCases = []testCase{
	{
		Name:  "Simple literal",
		Regex: "teapot",
		Want:  Literal("teapot"),
	},
	{
		Name:  "Empty literal",
		Regex: "",
		Want:  Sequence{},
	},
	{
		Name:  "Alternation",
		Regex: "teapot|flask|glass",
		Want:  Alternation{Literal("teapot"), Literal("flask"), Literal("glass")},
	},
	{
		Name:  "Empty alternation",
		Regex: "|",
		Want:  Alternation{Sequence{}, Sequence{}},
	},
	{
		Name:  "One-sided alternation (left)",
		Regex: "tea|",
		Want:  Alternation{Literal("tea"), Sequence{}},
	},
	{
		Name:  "One-sided alternation (right)",
		Regex: "|coffee",
		Want:  Alternation{Sequence{}, Literal("coffee")},
	},
	{
		Name:  "Grouping",
		Regex: "a(bc)",
		Want:  Sequence{Literal("a"), Group{Literal("bc")}},
	},
	{
		Name:  "Alternation with grouping",
		Regex: "b(a|o)ss",
		Want: Sequence{
			Literal("b"),
			Group{Alternation{Literal("a"), Literal("o")}},
			Literal("ss"),
		},
	},
	{
		Name:  "Repetition (*)",
		Regex: "yes*",
		Want:  Sequence{Literal("ye"), Repetition{Content: Literal("s"), UpperLimit: -1}},
	},
	{
		Name:  "Repetition (+)",
		Regex: "yes+",
		Want:  Sequence{Literal("ye"), Repetition{Content: Literal("s"), LowerLimit: 1, UpperLimit: -1}},
	},
	{
		Name:  "Repetition (counted, lower limit only)",
		Regex: "lo{2,}ng",
		Want:  Sequence{Literal("l"), Repetition{Content: Literal("o"), LowerLimit: 2, UpperLimit: -1}},
	},
	{
		Name:  "Repetition (counted, upper limit only)",
		Regex: "A{,3}",
		Want:  Repetition{Content: Literal("A"), UpperLimit: 3},
	},
	{
		Name:  "Repetition (counted, both limits)",
		Regex: "A{3,33}",
		Want:  Repetition{Content: Literal("A"), LowerLimit: 3, UpperLimit: 33},
	},
}

func TestParser(t *testing.T) {
	for _, re := range testCases {
		result, err := Parse(re.Regex)
		if err != nil {
			t.Errorf("%s: got %v, want %#v", re.Name, err, re.Want)
		}
		if !reflect.DeepEqual(result, re.Want) {
			t.Errorf("%s: got %#v, want %#v", re.Name, result, re.Want)
		}
	}
}
