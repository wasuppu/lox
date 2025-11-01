package interpreter

import (
	"fmt"
	"lox/treewalk/ast"
	"lox/treewalk/astprinter"
	"lox/treewalk/env"
	"lox/treewalk/loxerrors"
	"lox/treewalk/token"
	"reflect"
)

type Return struct {
	value any
}

func (r Return) Error() string {
	return fmt.Sprintf("%v", r.value)
}

type Interpreter struct {
	loxerror    *loxerrors.LoxErrors
	environment *env.Environment
	globals     *env.Environment
}

func New(loxerror *loxerrors.LoxErrors) *Interpreter {
	globals := env.New(loxerror, nil)
	globals.Define("clock", Clock{})
	return &Interpreter{loxerror: loxerror, environment: globals, globals: globals}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) (any, error) {
	var value any

	for _, statement := range statements {
		value = i.execute(statement)
		if err, ok := value.(error); ok {
			return nil, err
		}
	}

	return value, nil
}

func (i *Interpreter) execute(stmt ast.Stmt) any {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environemt *env.Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()
	i.environment = environemt

	for _, statement := range statements {
		i.execute(statement)
		// value := i.execute(statement)
		// if _, ok := value.(error); ok {
		// 	i.environment = previous
		// 	return
		// }
	}

}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) any {
	i.executeBlock(stmt.Statements, env.New(i.loxerror, i.environment))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) any {
	var value any
	if stmt.Initializer != nil {
		value = i.evalute(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) any {
	if isTruthy(i.evalute(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) any {
	for isTruthy(i.evalute(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	return nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) any {
	value := i.evalute(stmt.Expression)
	return value
}

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) any {
	function := LoxFunction{stmt}
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) any {
	value := i.evalute(stmt.Expression)
	if err, ok := value.(error); ok {
		return err
	}
	fmt.Println(astprinter.Stringify(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) any {
	var value any
	fmt.Println("type is: ", reflect.TypeOf(stmt))
	if stmt.Value != nil {
		fmt.Println("type of value: ", reflect.TypeOf(stmt.Value))
		fmt.Printf("value of stmt value: %#v\n", stmt.Value)

		value = i.evalute(stmt.Value)
	}
	// fmt.Println("value is: ", value)
	panic(Return{value})
}

func (i *Interpreter) evalute(exp ast.Expr) any {
	return exp.Accept(i)
}

func (i *Interpreter) VisitAssignExpr(exp *ast.Assign) any {
	value := i.evalute(exp.Value)
	err := i.environment.Assign(exp.Name, value)
	if err != nil {
		return err
	}
	return value
}

func (i *Interpreter) VisitVariableExpr(exp *ast.Variable) any {
	value, err := i.environment.Get(exp.Name)
	if err != nil {
		return err
	}
	return value
}

func (i *Interpreter) VisitLiteralExpr(exp *ast.Literal) any {
	return exp.Value
}

func (i *Interpreter) VisitLogicalExpr(exp *ast.Logical) any {
	left := i.evalute(exp.Left)

	if exp.Operator.Typ == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evalute(exp.Right)
}

func (i *Interpreter) VisitGroupingExpr(exp *ast.Grouping) any {
	return i.evalute(exp.Expression)
}

func (i *Interpreter) VisitUnaryExpr(exp *ast.Unary) any {
	right := i.evalute(exp.Right)
	op := exp.Operator
	switch op.Typ {
	case token.BANG:
		return !isTruthy(right)
	case token.MINUS:
		if r, ok := right.(float64); ok {
			return -r
		} else {
			err := loxerrors.NewErrorRuntime(op, "Operand must be a number")
			i.loxerror.RuntimeError(err)
			return err
		}
	}

	err := loxerrors.NewErrorRuntime(op, "Unreachable")
	i.loxerror.RuntimeError(err)
	return err
}

func (i *Interpreter) VisitBinaryExpr(exp *ast.Binary) any {
	left := i.evalute(exp.Left)
	if err, ok := left.(error); ok {
		return err
	}
	// fmt.Println("type of binary right: ", reflect.TypeOf(exp.Right))
	// fmt.Printf("value of binary right: %#v\n", exp.Right)
	right := i.evalute(exp.Right)
	if err, ok := right.(error); ok {
		return err
	}
	op := exp.Operator

	switch op.Typ {
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	default:
		l, ok1 := left.(float64)
		r, ok2 := right.(float64)
		if !(ok1 && ok2) {
			if op.Typ == token.PLUS {
				l, ok1 := left.(string)
				r, ok2 := right.(string)
				if ok1 && ok2 {
					return l + r
				}
			}

			err := loxerrors.NewErrorRuntime(op, "Operands must be two numbers or two strings.")
			i.loxerror.RuntimeError(err)

			return err
		}

		switch op.Typ {
		case token.GREATER:
			return l > r
		case token.GREATER_EQUAL:
			return l >= r
		case token.LESS:
			return l < r
		case token.LESS_EQUAL:
			return l <= r
		case token.MINUS:
			fmt.Println("minus")
			fmt.Println("l is", l)
			fmt.Println("r is", r)
			return l - r
		case token.PLUS:
			fmt.Println("plus")
			fmt.Println("l is", l)
			fmt.Println("r is", r)
			return l + r
		case token.SLASH:
			return l / r
		case token.STAR:
			return l * r
		}
	}

	err := loxerrors.NewErrorRuntime(op, "Unreachable")
	i.loxerror.RuntimeError(err)
	return err
}

func (i *Interpreter) VisitCallExpr(exp *ast.Call) any {
	callee := i.evalute(exp.Callee)

	var arguments []any
	for _, argument := range exp.Arguments {
		arguments = append(arguments, i.evalute(argument))
	}

	function, ok := callee.(LoxCallable)
	if !ok {
		err := loxerrors.NewErrorRuntime(exp.Paren, "Can only call functions and classes.")
		i.loxerror.RuntimeError(err)
		return err
	}

	if want, got := function.arity(), len(arguments); want != got {
		err := loxerrors.NewErrorRuntime(exp.Paren, fmt.Sprintf("Expected %d arguments but got %d.", want, got))
		i.loxerror.RuntimeError(err)
		return err
	}

	return function.call(i, arguments)
}

func isTruthy(o any) bool {
	if o == nil {
		return false
	}

	if b, ok := o.(bool); ok {
		return b
	}

	return true
}

func isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}
