package astprinter

import (
	"bytes"
	"fmt"
	"lox/treewalk/ast"
	"lox/treewalk/token"
	"strings"
)

type ASTPrinter struct{}

func New() *ASTPrinter {
	return &ASTPrinter{}
}

func Stringify(o any) string {
	if o == nil {
		return "nil"
	}

	return fmt.Sprint(o)
}

func (a ASTPrinter) PrintStmts(statements []ast.Stmt) string {
	// fmt.Println(statements)
	str := ""
	for i := 0; i < len(statements); i++ {
		str += a.PrintStmt(statements[i])
		// str += fmt.Sprint(statements[i].Accept(a).(string))
		if i < len(statements)-1 {
			str += fmt.Sprintln()
		}
	}
	return str
}

func (a ASTPrinter) PrintStmt(e ast.Stmt) string {
	return e.Accept(a).(string)
}

func (a ASTPrinter) printBlockWithIdent(stmt *ast.Block, level int) string {
	str := strings.Repeat("  ", level) + "{\n"
	for i := 0; i < len(stmt.Statements); i++ {
		if block, ok := stmt.Statements[i].(*ast.Block); ok {
			str += a.printBlockWithIdent(block, level+1)
		} else {
			str += strings.Repeat("  ", level+1) + fmt.Sprintln(stmt.Statements[i].Accept(a).(string))
		}
	}
	str += strings.Repeat("  ", level) + "}"
	if level != 0 {
		str += fmt.Sprintln()
	}
	return str
}

func (a ASTPrinter) VisitBlockStmt(stmt *ast.Block) any {
	return a.printBlockWithIdent(stmt, 0)
}

func (a ASTPrinter) VisitVarStmt(stmt *ast.Var) any {
	str := "var " + fmt.Sprintf("%v", stmt.Name.Lexeme)
	if stmt.Initializer != nil {
		str += " = " + a.Print(stmt.Initializer)
	}
	str += ";"
	return str
}

func (a ASTPrinter) VisitExpressionStmt(stmt *ast.Expression) any {
	return a.Print(stmt.Expression)
}

func (a ASTPrinter) VisitFunctionStmt(stmt *ast.Function) any {
	str := "fun " + stmt.Name.Lexeme + "("
	lst := len(stmt.Params) - 1
	if lst >= 0 {
		for i := 0; i < lst; i++ {
			str += stmt.Params[i].Lexeme + ", "
		}
		str += stmt.Params[lst].Lexeme
	}
	str += ") {\n" + a.PrintStmts(stmt.Body) + "\n}"
	return str
}

func (a ASTPrinter) VisitPrintStmt(stmt *ast.Print) any {
	return "print " + a.Print(stmt.Expression) + ";"
}

func (a ASTPrinter) VisitReturnStmt(stmt *ast.Return) any {
	return "return " + a.Print(stmt.Value) + ";"
}

func (a ASTPrinter) VisitWhileStmt(stmt *ast.While) any {
	// fmt.Println(reflect.TypeOf(stmt.Condition))
	// fmt.Println(reflect.TypeOf(stmt.Body.(*ast.Print).Expression.(*ast.Assign).Value))
	return "while (" + a.Print(stmt.Condition) + ") " + a.PrintStmt(stmt.Body)
}

func (a ASTPrinter) VisitIfStmt(stmt *ast.If) any {
	str := "if (" + a.Print(stmt.Condition) + ") " + a.PrintStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		str += " else " + a.PrintStmt(stmt.ElseBranch)
	}
	return str
}

func (a ASTPrinter) VisitAssignExpr(e *ast.Assign) any {
	return fmt.Sprintf("%v", e.Name.Lexeme) + " = " + a.Print(e.Value)
}

func (a ASTPrinter) Print(e ast.Expr) string {
	return e.Accept(a).(string)
}

func (a ASTPrinter) VisitVariableExpr(e *ast.Variable) any {
	return fmt.Sprintf("%v", e.Name.Lexeme)
}

func (a ASTPrinter) VisitLiteralExpr(e *ast.Literal) any {
	if e.Value == nil {
		return "null"
	}

	if s, ok := e.Value.(string); ok {
		return fmt.Sprintf(`"%s"`, s)
	}

	return fmt.Sprintf("%v", e.Value)
}

func (a ASTPrinter) VisitLogicalExpr(e *ast.Logical) any {
	var op string
	if e.Operator.Typ == token.OR {
		op = " or "
	} else {
		op = " and "
	}
	return a.Print(e.Left) + op + a.Print(e.Right)
}

func (a ASTPrinter) VisitGroupingExpr(e *ast.Grouping) any {
	return a.parenthesize("", e.Expression)
	// return a.parenthesize("group", e.Expression)
}

func (a ASTPrinter) VisitUnaryExpr(e *ast.Unary) any {
	return a.parenthesize(e.Operator.Lexeme, e.Right)
}

func (a ASTPrinter) VisitBinaryExpr(e *ast.Binary) any {
	return a.Print(e.Left) + fmt.Sprintf(" %s ", e.Operator.Lexeme) + a.Print(e.Right)
	// return a.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (a ASTPrinter) VisitCallExpr(e *ast.Call) any {
	str := a.Print(e.Callee) + "("
	lst := len(e.Arguments) - 1
	if lst >= 0 {
		for i := 0; i < lst; i++ {
			str += a.Print(e.Arguments[i]) + ", "
		}
		str += a.Print(e.Arguments[lst])
	}
	str += ")"
	return str
}

func (a ASTPrinter) parenthesize(name string, exprs ...ast.Expr) string {
	buf := bytes.Buffer{}

	buf.WriteRune('(')
	buf.WriteString(name)
	for _, expr := range exprs {
		// buf.WriteRune(' ')
		v, ok := expr.Accept(a).(string)
		if !ok {
			v = ""
		}
		buf.WriteString(v)
	}
	buf.WriteRune(')')
	return buf.String()
}
