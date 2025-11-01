package astprinter

import (
	"fmt"
	"lox/treewalk/ast"
	"lox/treewalk/token"
	"testing"
)

func TestPrinter(t *testing.T) {
	printer := ASTPrinter{}
	var exp ast.Expr = ast.NewBinary(
		ast.NewUnary(token.New(token.MINUS, "-", nil, 1), ast.NewLiteral(123)),
		token.New(token.STAR, "*", nil, 1),
		ast.NewGrouping(ast.NewLiteral(45.67)))
	fmt.Println(printer.Print(exp))
}
