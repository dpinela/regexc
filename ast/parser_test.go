package ast

import (
	"reflect"
	"testing"
)

type testCase struct {
	Name  string
	Regex string
	Want  interface{}
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
		Want:  Sequence{Literal("l"), Repetition{Content: Literal("o"), LowerLimit: 2, UpperLimit: -1}, Literal("ng")},
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
	{
		Name:  "Repetition (counted, fixed)",
		Regex: "($){12}",
		Want:  Repetition{Content: Group{Literal("$")}, LowerLimit: 12, UpperLimit: 12},
	},
	{
		Name:  "Empty group",
		Regex: "(())",
		Want:  Group{Group{Sequence{}}},
	},
	{
		Name:  "Unterminated group",
		Regex: "(endless",
		Want:  &UnterminatedGroupError{Location: 8, Source: "(endless"},
	},
	{
		Name:  "Extra closing parenthesis",
		Regex: "(green) )tea",
		Want:  &BadCloseError{Location: 8, Source: "(green) )tea"},
	},
	{
		Name:  "Repetition of nothing",
		Regex: "*",
		Want:  &VoidRepetitionError{Location: 0, Source: "*"},
	},
	{
		Name:  "Repetition of nothing 2",
		Regex: "({2,5})",
		Want:  &VoidRepetitionError{Location: 1, Source: "({2,5})"},
	},
	{
		Name:  "Repetition of repetition",
		Regex: "bo++m",
		Want:  &RepetitionRepetitionError{Location: 3, Source: "bo++m"},
	},
	{
		Name:  "Invalid counted repetition (2 commas)",
		Regex: "(x){2,3,2}",
		Want:  &RepetitionBadCharError{parseError: parseError{Location: 7, Source: "(x){2,3,2}"}, Char: ','},
	},
	{
		Name:  "Invalid counted repetition (invalid char)",
		Regex: "($){$,$}",
		Want:  &RepetitionBadCharError{parseError: parseError{Location: 4, Source: "($){$,$}"}, Char: '$'},
	},
	{
		Name:  "Closing non-existent counted repetition",
		Regex: "}{",
		Want:  &BadRepetitionCloseError{Location: 0, Source: "}{"},
	},
	{
		Name:  "Unterminated counted repetition",
		Regex: "(forever){4,",
		Want:  &UnterminatedRepetitionError{Location: 12, Source: "(forever){4,"},
	},
}

// resultOfParsing returns the tree obtained by parsing re, if it is valid, or the error
// otherwise.
func resultOfParsing(re string) interface{} {
	tree, err := Parse(re)
	if err != nil {
		return err
	}
	return tree
}

func TestParser(t *testing.T) {
	for _, re := range testCases {
		if result := resultOfParsing(re.Regex); !reflect.DeepEqual(result, re.Want) {
			t.Errorf("%s: got %#v, want %#v", re.Name, result, re.Want)
		}
	}
}
