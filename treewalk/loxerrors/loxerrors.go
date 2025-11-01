package loxerrors

import (
	"errors"
	"fmt"
	"lox/treewalk/token"
	"os"
)

var (
	ErrorParse = errors.New("ParseError")
)

type ErrorRuntime struct {
	token   token.Token
	message string
}

func (r ErrorRuntime) Error() string {
	return r.message
}

func NewErrorRuntime(token token.Token, message string) *ErrorRuntime {
	return &ErrorRuntime{token: token, message: message}
}

type LoxErrors struct {
	HadError        bool
	HadRuntimeError bool
}

func New() *LoxErrors {
	return &LoxErrors{}
}

func (le *LoxErrors) RuntimeError(err *ErrorRuntime) {
	fmt.Fprintf(os.Stderr, err.message+"\n[line %d]\n", err.token.Line)
	le.HadError = true
}

func (le *LoxErrors) TokenError(tok token.Token, message string) {
	if tok.Typ == token.EOF {
		le.report(tok.Line, " at end", message)
	} else {
		le.report(tok.Line, " at '"+tok.Lexeme+"'", message)
	}
}

func (le *LoxErrors) Error(line int, message string) {
	le.report(line, "", message)
}

func (le *LoxErrors) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)
	le.HadError = true
}
