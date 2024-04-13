package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"monkeylang/token"
	"strconv"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	LOWEST int = iota
	EQUALITY
	COMPARISON
	TERM
	FACTOR
	UNARY
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUAL:        EQUALITY,
	token.BANG_EQUAL:   EQUALITY,
	token.LESS_THAN:    COMPARISON,
	token.GREATER_THAN: COMPARISON,
	token.PLUS:         TERM,
	token.MINUS:        TERM,
	token.SLASH:        FACTOR,
	token.ASTERISK:     FACTOR,
}

type Parser struct {
	l *lexer.Lexer

	Errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		Errors: []string(nil),
	}

	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixParseFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixParseFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParseFn(token.BANG, p.parseUnaryExpression)
	p.registerPrefixParseFn(token.MINUS, p.parseUnaryExpression)

    p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixParseFn(token.EQUAL, p.parseBinaryExpression)
	p.registerInfixParseFn(token.BANG_EQUAL, p.parseBinaryExpression)
	p.registerInfixParseFn(token.LESS_THAN, p.parseBinaryExpression)
	p.registerInfixParseFn(token.GREATER_THAN, p.parseBinaryExpression)
	p.registerInfixParseFn(token.PLUS, p.parseBinaryExpression)
	p.registerInfixParseFn(token.MINUS, p.parseBinaryExpression)
	p.registerInfixParseFn(token.SLASH, p.parseBinaryExpression)
	p.registerInfixParseFn(token.ASTERISK, p.parseBinaryExpression)

	return p
}

func (p *Parser) registerPrefixParseFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixParseFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) reportMatchError(t token.TokenType) {
	msg := fmt.Sprintf(
		"expected next token to be %s, got %s instead",
		t, p.peekToken.Type,
	)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement(nil)

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.matchNext(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

	if !p.matchNext(token.ASSIGN) {
		return nil
	}

	// TODO: skipping expressions
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.IdentifierExpression{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.Errors = append(p.Errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parseUnaryExpression() ast.Expression {
	expr := &ast.UnaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expr.Right = p.parseExpression(UNARY)

	return expr
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: skipping expressions
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if !p.matchNext(token.SEMICOLON) {
		p.reportMatchError(token.SEMICOLON)
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.reportMissingPrefixParseFn(p.curToken.Type)
		return nil
	}

	left := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return left
		}

		p.nextToken()
        left = p.parseBinaryExpression(left)
	}

	return left
}

func (p *Parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

    precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)

	return expr
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedences[p.curToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}
func (p *Parser) reportMissingPrefixParseFn(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse fn for %s found", t)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) matchNext(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.reportMatchError(t)
	return false
}
