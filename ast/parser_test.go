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
		Regex: "teapot|flask",
		Want:  Alternation{Literal("teapot"), Literal("flask")},
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
}

func TestParser(t *testing.T) {
	for _, re := range testCases {
		if result, err := Parse(re.Regex); err != nil || !reflect.DeepEqual(result, re.Want) {
			t.Errorf("%s: got %#v, want %#v", re.Name, result, re.Want)
		}
	}
}
