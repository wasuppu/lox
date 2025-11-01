package ast

import "lox/treewalk/token"

type StmtVisitor interface {
	VisitPrintStmt(expr *Print) any
	VisitReturnStmt(expr *Return) any
	VisitVarStmt(expr *Var) any
	VisitBlockStmt(expr *Block) any
	VisitExpressionStmt(expr *Expression) any
	VisitFunctionStmt(expr *Function) any
	VisitIfStmt(expr *If) any
	VisitWhileStmt(expr *While) any
}

type Stmt interface {
	Accept(v StmtVisitor) any
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Stmt {
	return &Print{Expression: expression}
}

func (e *Print) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(e)
}

type Return struct {
	Keyword token.Token
	Value   Expr
}

func NewReturn(keyword token.Token, value Expr) Stmt {
	return &Return{Keyword: keyword, Value: value}
}

func (e *Return) Accept(v StmtVisitor) any {
	return v.VisitReturnStmt(e)
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

func NewVar(name token.Token, initializer Expr) Stmt {
	return &Var{Name: name, Initializer: initializer}
}

func (e *Var) Accept(v StmtVisitor) any {
	return v.VisitVarStmt(e)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) Stmt {
	return &Block{Statements: statements}
}

func (e *Block) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(e)
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{Expression: expression}
}

func (e *Expression) Accept(v StmtVisitor) any {
	return v.VisitExpressionStmt(e)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func NewFunction(name token.Token, params []token.Token, body []Stmt) Stmt {
	return &Function{Name: name, Params: params, Body: body}
}

func (e *Function) Accept(v StmtVisitor) any {
	return v.VisitFunctionStmt(e)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenbranch Stmt, elsebranch Stmt) Stmt {
	return &If{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}

func (e *If) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(e)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return &While{Condition: condition, Body: body}
}

func (e *While) Accept(v StmtVisitor) any {
	return v.VisitWhileStmt(e)
}
