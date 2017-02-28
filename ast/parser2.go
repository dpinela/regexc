package ast

type Node interface {
	simplify() Node
}

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

func Parse(re string) (Node, error) {
	p := parser{stack: []Node{Sequence(nil)}}
	tree, err := p.parseRegexp(re)
	if err != nil {
		return nil, err
	}
	return tree.simplify(), nil
}

type parser struct {
	stack []Node
}

type parseError struct {
	Message  string
	Location int
	Source   string
}

type BadCloseError parseError

func (err *BadCloseError) Error() string { return "closing parenthesis outside of group" }

type UnterminatedGroupError parseError

func (err *UnterminatedGroupError) Error() string { return "unterminated group" }

type VoidRepetitionError parseError

func (err *VoidRepetitionError) Error() string { return "repetition of nothing" }

type RepetitionRepetitionError parseError

func (err *RepetitionRepetitionError) Error() string { return "repetition of repetition" }

func (p *parser) parseRegexp(re string) (Node, error) {
	//fmt.Println("Parsing", re)
	groupLevel := 0
	for i, c := range re {
		switch c {
		case '*', '+':
		default:
			p.extendSequence()
		}
		//fmt.Printf("char: %c stack: %#v\n", c, p.stack)
		switch c {
		case '(':
			// Add a group token to the stack so that combining operations don't mix
			// the elements of the group with elements of the group's parent
			// (like extendSequence)
			p.stack = append(p.stack, Group{}, Sequence{})
			groupLevel++
		case ')':
			p.extendAlternation(true)
			if !p.finishGroup() {
				return nil, &BadCloseError{Location: i, Source: re}
			}
			groupLevel--
		case '|':
			p.startOrExtendAlternation()
		case '*':
			if err := p.addRepetition(0, -1, i, re); err != nil {
				return nil, err
			}
		case '+':
			if err := p.addRepetition(1, -1, i, re); err != nil {
				return nil, err
			}
		default:
			p.stack = append(p.stack, Literal(c))
		}
	}
	p.extendSequence()
	p.extendAlternation(true)
	//fmt.Printf("finish, stack: %#v\n", p.stack)
	if groupLevel > 0 {
		return nil, &UnterminatedGroupError{Location: len(re), Source: re}
	}
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
		case Group:
			p.push(target)
			if finish {
				p.push(bit)
			} else {
				p.push(Alternation{bit})
				p.push(Sequence{})
			}
		default:
			p.push(target)
			p.push(bit)
		}
	}
}

func (p *parser) addRepetition(lowerLimit, upperLimit, pos int, re string) error {
	if len(p.stack) == 0 {
		return &VoidRepetitionError{Location: pos, Source: re}
	}
	switch target := p.pop().(type) {
	case Sequence: // can happen when the repeat operator appears at the start of a sequence
		return &VoidRepetitionError{Location: pos, Source: re}
	case Alternation:
		panic("BUG: addRepetition called with an Alternation at top of stack")
	case Repetition:
		return &RepetitionRepetitionError{Location: pos, Source: re}
	default:
		p.push(Repetition{Content: target, LowerLimit: lowerLimit, UpperLimit: upperLimit})
		return nil
	}
}

func (p *parser) finishGroup() bool {
	if len(p.stack) < 2 {
		return false
	}
	content := p.pop()
	if _, isGroup := p.pop().(Group); isGroup {
		p.push(Group{content})
	} else {
		return false
	}
	return true
}

func (p *parser) push(item Node) {
	p.stack = append(p.stack, item)
}

func (p *parser) pop() Node {
	item := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return item
}
