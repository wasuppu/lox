package env

import (
	"lox/treewalk/loxerrors"
	"lox/treewalk/token"
)

type Environment struct {
	Values    map[string]any
	enclosing *Environment
	loxerror  *loxerrors.LoxErrors
}

func New(loxerror *loxerrors.LoxErrors, enclosing *Environment) *Environment {
	env := Environment{loxerror: loxerror, enclosing: enclosing}
	env.Values = make(map[string]any)
	return &env
}

func Copy(e1 *Environment) *Environment {
	values := make(map[string]any)
	for k, v := range e1.Values {
		values[k] = v
	}
	var enclosing *Environment
	if e1.enclosing != nil {
		enclosing = Copy(e1.enclosing)
	}
	return &Environment{values, enclosing, e1.loxerror}
}

func (env *Environment) Get(name token.Token) (any, error) {
	if value, ok := env.Values[name.Lexeme]; ok {
		return value, nil
	}

	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}

	err := loxerrors.NewErrorRuntime(name, "Undefined variable '"+name.Lexeme+"'.")
	env.loxerror.RuntimeError(err)
	return nil, err
}

func (env *Environment) Define(name string, value any) {
	env.Values[name] = value
}

func (env *Environment) Assign(name token.Token, value any) error {
	if _, ok := env.Values[name.Lexeme]; ok {
		env.Values[name.Lexeme] = value
		return nil
	}

	if env.enclosing != nil {
		env.enclosing.Assign(name, value)
		return nil
	}

	err := loxerrors.NewErrorRuntime(name, "Undefined variable '"+name.Lexeme+"'.")
	env.loxerror.RuntimeError(err)
	return err
}
