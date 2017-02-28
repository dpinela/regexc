package ast

import "fmt"

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

func Parse(re string) Node {
	p := parser{stack: []Node{Sequence(nil)}}
	tree, _ := p.parseRegexp(re)
	return tree
}

type parser struct {
	stack []Node
}

func (p *parser) parseRegexp(re string) (Node, error) {
	fmt.Println("Parsing", re)
	for _, c := range re {
		fmt.Printf("char: %c stack: %#v\n", c, p.stack)
		p.extendSequence()
		switch c {
		case '|':
			p.startOrExtendAlternation()
		default:
			p.stack = append(p.stack, Literal(c))
		}
	}
	p.extendSequence()
	p.extendAlternation(true)
	return p.pop(), nil
}

func (p *parser) extendSequence() {
	if len(p.stack) < 2 {
		return
	}
	bit := p.pop()
	switch target := p.pop().(type) {
	case Sequence:
		p.push(append(target, bit))
	default:
		p.push(target)
		p.push(bit)
	}
}

// on seeing |
// pop
func (p *parser) startOrExtendAlternation() {
	switch len(p.stack) {
	case 0:
	case 1:
		p.push(Alternation{p.pop()})
		p.push(Sequence{})
	default:
		p.extendAlternation(false)
	}
}

func (p *parser) extendAlternation(finish bool) {
	if len(p.stack) >= 2 {
		bit := p.pop()
		switch target := p.pop().(type) {
		case Alternation:
			p.push(append(target, bit))
			if !finish {
				p.push(Sequence{})
			}
		default:
			p.push(target)
			p.push(bit)
		}
	}
}

func (p *parser) push(item Node) {
	p.stack = append(p.stack, item)
}

func (p *parser) pop() Node {
	item := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return item
}
