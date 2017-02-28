package ast

func (l Literal) simplify() Node     { return l }
func (g Group) simplify() Node       { return Group{g.Content.simplify()} }
func (a Alternation) simplify() Node { return Alternation(simplifySubtree([]Node(a))) }
func (s Sequence) simplify() Node {
	joinedS := consolidateLiteralRuns([]Node(s))
	if len(joinedS) == 1 {
		return joinedS[0].simplify()
	}
	return Sequence(simplifySubtree([]Node(joinedS)))
}

func (r Repetition) simplify() Node {
	return Repetition{Content: r.Content.simplify(), LowerLimit: r.LowerLimit, UpperLimit: r.UpperLimit}
}

func consolidateLiteralRuns(ns []Node) []Node {
	joinedLit := Literal("")
	var newNs []Node
	for _, n := range ns {
		if lit, ok := n.(Literal); ok {
			joinedLit += lit
		} else {
			if joinedLit != "" {
				newNs = append(newNs, joinedLit)
			}
			newNs = append(newNs, n)
			joinedLit = ""
		}
	}
	if joinedLit != "" {
		newNs = append(newNs, joinedLit)
	}
	return newNs
}

func simplifySubtree(ns []Node) []Node {
	newNs := make([]Node, len(ns))
	for i, n := range ns {
		newNs[i] = n.simplify()
	}
	return newNs
}
