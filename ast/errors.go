package ast

import "strconv"

type parseError struct {
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

type RepetitionBadCharError struct {
	parseError
	Char rune
}

func (err *RepetitionBadCharError) Error() string {
	return "unexpected " + strconv.QuoteRune(err.Char) + " in counted repetition"
}

type UnterminatedRepetitionError parseError

func (err *UnterminatedRepetitionError) Error() string { return "unterminated counted repetition" }

type UnterminatedCharClassError parseError

func (err *UnterminatedCharClassError) Error() string { return "unterminated character class" }

type BadRepetitionCloseError parseError

func (err *BadRepetitionCloseError) Error() string {
	return "closing brace outside of counted repetition"
}

type BadCharClassCloseError parseError

func (err *BadCharClassCloseError) Error() string { return "closing bracket outside of character class" }
