package ast

import "lox/treewalk/token"

type ExprVisitor interface {
	VisitLiteralExpr(expr *Literal) any
	VisitGroupingExpr(expr *Grouping) any
	VisitUnaryExpr(expr *Unary) any
	VisitLogicalExpr(expr *Logical) any
	VisitBinaryExpr(expr *Binary) any
	VisitCallExpr(expr *Call) any
	VisitVariableExpr(expr *Variable) any
	VisitAssignExpr(expr *Assign) any
}

type Expr interface {
	Accept(v ExprVisitor) any
}

type Literal struct {
	Value any
}

func NewLiteral(value any) Expr {
	return &Literal{Value: value}
}

func (e *Literal) Accept(v ExprVisitor) any {
	return v.VisitLiteralExpr(e)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{Expression: expression}
}

func (e *Grouping) Accept(v ExprVisitor) any {
	return v.VisitGroupingExpr(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func NewUnary(operator token.Token, right Expr) Expr {
	return &Unary{Operator: operator, Right: right}
}

func (e *Unary) Accept(v ExprVisitor) any {
	return v.VisitUnaryExpr(e)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogical(left Expr, operator token.Token, right Expr) Expr {
	return &Logical{Left: left, Operator: operator, Right: right}
}

func (e *Logical) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(e)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) Expr {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func (e *Binary) Accept(v ExprVisitor) any {
	return v.VisitBinaryExpr(e)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren token.Token, arguments []Expr) Expr {
	return &Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (e *Call) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(e)
}

type Variable struct {
	Name token.Token
}

func NewVariable(name token.Token) Expr {
	return &Variable{Name: name}
}

func (e *Variable) Accept(v ExprVisitor) any {
	return v.VisitVariableExpr(e)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func NewAssign(name token.Token, value Expr) Expr {
	return &Assign{Name: name, Value: value}
}

func (e *Assign) Accept(v ExprVisitor) any {
	return v.VisitAssignExpr(e)
}
