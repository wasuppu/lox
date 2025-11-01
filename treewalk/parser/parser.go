package parser

import (
	"lox/treewalk/ast"
	"lox/treewalk/loxerrors"
	"lox/treewalk/token"
)

type Parser struct {
	tokens   []token.Token
	current  int
	loxerror *loxerrors.LoxErrors
}

func New(tokens []token.Token, loxerror *loxerrors.LoxErrors) *Parser {
	return &Parser{tokens: tokens, loxerror: loxerror}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	var err error

	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			p.synchronize()
		}
		statements = append(statements, statement)
	}

	return statements, err
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(token.FUN) {
		return p.function("function")
	}

	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.FOR) {
		return p.forStatement()
	}

	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.RETURN) {
		return p.returnStatement()
	}

	if p.match(token.WHILE) {
		return p.whileStatement()
	}

	if p.match(token.LEFT_BRACE) {
		block, err := p.block()
		return ast.NewBlock(block), err
	}

	return p.expressionStatement()
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	return ast.NewVar(name, initializer), nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	var err error
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var contidition ast.Expr
	if !p.check(token.SEMICOLON) {
		contidition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = ast.NewBlock(append([]ast.Stmt{body}, ast.NewExpression(increment)))
	}
	if contidition == nil {
		contidition = ast.NewLiteral(true)
	}
	body = ast.NewWhile(contidition, body)
	if initializer != nil {
		body = ast.NewBlock(append([]ast.Stmt{initializer}, body))
	}

	return body, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return ast.NewWhile(condition, body), nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	var statements []ast.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return ast.NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return ast.NewPrint(value), nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	var err error
	keyword := p.previous()
	var value ast.Expr
	if !p.check(token.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(token.SEMICOLON, "Expect ';' after return value.")
	return ast.NewReturn(keyword, value), nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	exp, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return ast.NewExpression(exp), nil
}

func (p *Parser) function(kind string) (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect "+kind+"name.")
	if err != nil {
		return nil, err
	}
	p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	var parameters []token.Token
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.loxerror.TokenError(p.peek(), "Can't have more than 255 parameters.")
			}

			id, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, id)

			if !p.match(token.COMMA) {
				break
			}
		}
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")

	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return ast.NewFunction(name, parameters, body), nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	exp, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := exp.(*ast.Variable); ok {
			return ast.NewAssign(variable.Name, value), nil
		}

		p.loxerror.TokenError(equals, "Invalid assignment target.")
	}

	return exp, nil
}

func (p *Parser) or() (ast.Expr, error) {
	exp, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.match(token.OR) {
		op := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		exp = ast.NewLogical(exp, op, right)
	}

	return exp, nil
}

func (p *Parser) and() (ast.Expr, error) {
	exp, err := p.equality()
	if err != nil {
		return nil, err
	}
	for p.match(token.AND) {
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		exp = ast.NewLogical(exp, op, right)
	}

	return exp, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	exp, err := p.comparision()
	if err != nil {
		return nil, err
	}
	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		op := p.previous()
		right, err := p.comparision()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, op, right)
	}

	return exp, nil
}

func (p *Parser) comparision() (ast.Expr, error) {
	exp, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, op, right)
	}

	return exp, nil
}

func (p *Parser) term() (ast.Expr, error) {
	exp, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.MINUS, token.PLUS) {
		op := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, op, right)
	}

	return exp, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	exp, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		exp = ast.NewBinary(exp, op, right)
	}
	return exp, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		op := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return ast.NewUnary(op, right), nil
	}

	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {
	exp, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			exp, err = p.finishCall(exp)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return exp, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	var arguments []ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		for {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			if len(arguments) >= 255 {
				p.loxerror.TokenError(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, exp)
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return ast.NewCall(callee, paren, arguments), nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return ast.NewLiteral(false), nil
	}
	if p.match(token.TRUE) {
		return ast.NewLiteral(true), nil
	}
	if p.match(token.NIL) {
		return ast.NewLiteral(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.NewLiteral(p.previous().Literal), nil
	}

	if p.match(token.LEFT_PAREN) {
		exp, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return ast.NewGrouping(exp), nil
	}

	if p.match(token.IDENTIFIER) {
		return ast.NewVariable(p.previous()), nil
	}

	p.loxerror.TokenError(p.peek(), "Expect expression.")
	return nil, loxerrors.ErrorParse
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p Parser) check(typ token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Typ == typ
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().Typ == token.EOF
}

func (p Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(typ token.TokenType, message string) (token.Token, error) {
	if p.check(typ) {
		return p.advance(), nil
	}
	p.loxerror.TokenError(p.peek(), message)
	return token.Token{}, loxerrors.ErrorParse
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Typ == token.SEMICOLON {
			return
		}

		switch p.peek().Typ {
		case token.CLASS, token.FUN, token.VAR, token.FOR:
			return
		case token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
