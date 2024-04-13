package ast

import (
	"bytes"
	"monkeylang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for i := range len(p.Statements) {
		out.WriteString(p.Statements[i].String())
	}

	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token
	Name  *IdentifierExpression
	Value Expression
}

func (stmt *LetStatement) statementNode()       {}
func (stmt *LetStatement) TokenLiteral() string { return stmt.Token.Literal }

func (stmt *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(stmt.TokenLiteral() + " ")
	out.WriteString(stmt.Name.String())
	out.WriteString(" = ")
	if stmt.Value != nil {
		out.WriteString(stmt.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type IdentifierExpression struct {
	Token token.Token
	Value string
}

func (expr *IdentifierExpression) expressionNode()      {}
func (expr *IdentifierExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *IdentifierExpression) String() string       { return expr.Value }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (stmt *ReturnStatement) statementNode()       {}
func (stmt *ReturnStatement) TokenLiteral() string { return stmt.Token.Literal }

func (stmt *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(stmt.TokenLiteral() + " ")
	if stmt.ReturnValue != nil {
		out.WriteString(stmt.ReturnValue.String())
	}
	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (stmt *ExpressionStatement) statementNode()       {}
func (stmt *ExpressionStatement) TokenLiteral() string { return stmt.Token.Literal }

func (stmt *ExpressionStatement) String() string {
	if stmt.Expression != nil {
		return stmt.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (expr *IntegerLiteral) expressionNode()      {}
func (expr *IntegerLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *IntegerLiteral) String() string       { return expr.Token.Literal }

type UnaryExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (expr *UnaryExpression) expressionNode()      {}
func (expr *UnaryExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *UnaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Operator)
	out.WriteString(expr.Right.String())
	out.WriteString(")")

	return out.String()
}
