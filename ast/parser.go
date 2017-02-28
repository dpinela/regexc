//+build ignore

package ast

import (
	"fmt"
	"regexp"
	"strings"
)

type Node interface{}

type Literal string

type Repetition struct {
	Content                Node
	LowerLimit, UpperLimit int
}

type Group struct {
	Content Node
}

type Sequence []Node

type Alternation []Node

type CharClass struct {
	Negated bool
	Ranges  []CharRange
}

type CharRange struct {
	Min, Max rune
}

// Regexp := Sequence | Alternation
// Sequence := RegexpElement*
// RegexpElement := Literal | CharClass | Repetition | Group
// Alternation := Sequence ('|' Sequence)+
// Literal := LiteralChar+
// CharClass := '[' CharRange+ ']'
// CharRange := EscAny | ([^-\]] | EscAny) '-' ([^-\]] | EscAny)
// Repetition := (LiteralChar | CharClass | Group) ([*+] | '{' \d* ',' \d* '}')
// Group := '(' Regexp ')'

// LiteralChar := EscAny | [^?*.\[\]()|]
// EscAny := '\\' Any
func Parse(re string) Node {
	fmt.Println("Parsing", re)
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(abortParse); !ok {
				panic(e)
			}
		}
	}()
	p := parser{text: re}
	p.parseRegexp(func() {})
	if p.pos < len(p.text) || len(p.resultStack) == 0 {
		return nil
	}
	return p.resultStack[0]
}

type parser struct {
	stack       []func()
	resultStack []interface{}
	pos         int
	text        string
}

type abortParse struct{}

func (p *parser) trace(s string) {
	fmt.Println(strings.Repeat(" ", len(p.stack))+s, "@", p.pos)
}

func (p *parser) Backtrack() {
	if len(p.stack) == 0 {
		panic(abortParse{})
	}
	p.trace("Backtrack")
	p.pop()()
}

func (p *parser) pop() func() {
	n := len(p.stack) - 1
	next := p.stack[n]
	p.stack = p.stack[:n]
	return next
}

func (p *parser) popResult() interface{} {
	n := len(p.resultStack) - 1
	res := p.resultStack[n]
	p.resultStack = p.resultStack[:n]
	return res
}

func dupEfaceSlice(xs []interface{}) []interface{} {
	xs2 := make([]interface{}, len(xs))
	copy(xs2, xs)
	return xs2
}

func (p *parser) Choose(fs []func(func()), then func()) {
	if len(fs) == 0 {
		p.Backtrack()
		return
	}
	savedPos := p.pos
	savedStack := dupEfaceSlice(p.resultStack)
	p.stack = append(p.stack, func() {
		p.pos = savedPos
		p.resultStack = savedStack
		p.Choose(fs[1:], then)
	})
	fs[0](then)
}

func (p *parser) matchPattern(re *regexp.Regexp, then func([]string)) {
	m := re.FindStringSubmatch(p.text[p.pos:])
	if len(m) == 0 {
		p.Backtrack()
		return
	}
	p.pos += len(m[0])
	then(m)
}

func (p *parser) matchByte(b byte, then func()) {
	if p.pos >= len(p.text) || p.text[p.pos] != b {
		p.Backtrack()
		return
	}
	p.pos++
	then()
}

func nop(then func()) { then() }

func (p *parser) parseRegexp(then func()) {
	p.trace("parseRegexp")
	p.Choose([]func(func()){p.parseAlternation, p.parseSequence}, then)
}

func (p *parser) parseSequence(then func()) {
	p.trace("parseSequence")
	p.parseOneOrMore(p.parseRegexpElement, nil,
		func(elements []Node) Node {
			if len(elements) == 1 {
				return elements[0]
			}
			return Sequence(elements)
		}, then)
}

func (p *parser) parseOneOrMore(parseFunc func(func()), elements []Node, typeConverter func([]Node) Node, then func()) {
	parseFunc(func() {
		p.Choose([]func(func()){
			func(then func()) {
				p.parseOneOrMore(parseFunc, append(elements, p.popResult()), typeConverter, then)
			}, func(then func()) {
				p.resultStack = append(p.resultStack, typeConverter(append(elements, p.popResult())))
				then()
			},
		}, then)
	})
}

func (p *parser) parseZeroOrMore(parseFunc func(func()), elements []Node, typeConverter func([]Node) Node, then func()) {
	p.Choose([]func(func()){
		func(then func()) { p.parseOneOrMore(parseFunc, elements, typeConverter, then) },
		func(then func()) {
			p.resultStack = append(p.resultStack, typeConverter(nil))
		},
	}, then)
}

func (p *parser) parseRegexpElement(then func()) {
	p.trace("parseRegexpElement")
	p.Choose([]func(func()){p.parseLiteral, p.parseGroup}, then)
}

var literalRE = regexp.MustCompile(`^(\.|[^?*.\[\]()|])+`)

func (p *parser) parseLiteral(then func()) {
	p.trace("parseLiteral")
	p.matchPattern(literalRE, func(match []string) {
		p.resultStack = append(p.resultStack, Literal(match[0]))
		then()
	})
}

func (p *parser) parseGroup(then func()) {
	p.trace("parseGroup")
	p.matchByte('(', func() {
		p.parseRegexp(func() {
			p.matchByte(')', func() {
				n := len(p.resultStack) - 1
				p.resultStack[n] = Group{Content: p.resultStack[n]}
				then()
			})
		})
	})
}

func (p *parser) parseCharClass(func())  {}
func (p *parser) parseRepetition(func()) {}

func (p *parser) parseAlternation(then func()) {
	fmt.Println("parseAlternation @", p.pos)
	p.parseSequence(func() {
		p.parseOneOrMore(func(then func()) {
			p.matchByte('|', func() { p.parseSequence(then) })
		}, []Node{p.popResult()}, func(ns []Node) Node { return Alternation(ns) }, then)
	})
}
