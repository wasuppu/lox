package interpreter

import (
	"fmt"
	"lox/treewalk/ast"
	"lox/treewalk/env"
	"time"
)

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) any
}

type LoxFunction struct {
	declaration *ast.Function
}

func (fn LoxFunction) call(interpreter *Interpreter, arguments []any) (returnValue any) {
	// globals := env.Copy(interpreter.globals)
	environment := env.New(interpreter.loxerror, interpreter.globals)
	for i, param := range fn.declaration.Params {
		environment.Define(param.Lexeme, arguments[i])
	}

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(Return); ok {
				fmt.Println(v)
				returnValue = v.value
				return
			} else {
				panic(err)
			}
		}
	}()

	interpreter.executeBlock(fn.declaration.Body, environment)
	return returnValue
}

func (fn LoxFunction) arity() int {
	return len(fn.declaration.Params)
}

func (fn LoxFunction) String() string {
	return "<fn " + fn.declaration.Name.Lexeme + ">"
}

type Clock struct{}

func (Clock) arity() int {
	return 0
}

func (Clock) call(interpreter *Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli() / 1000.0)
}

func (Clock) String() string {
	return "<native fn Clock>"
}
