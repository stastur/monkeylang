package ast

import (
	"bytes"
	"monkeylang/token"
	"strings"
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

type BinaryExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

func (expr *BinaryExpression) expressionNode()      {}
func (expr *BinaryExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *BinaryExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(expr.Left.String())
	out.WriteString(" " + expr.Operator + " ")
	out.WriteString(expr.Right.String())
	out.WriteString(")")

	return out.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (expr *Boolean) expressionNode()      {}
func (expr *Boolean) TokenLiteral() string { return expr.Token.Literal }
func (expr *Boolean) String() string       { return expr.Token.Literal }

type IfExpression struct {
	Token      token.Token
	Condition  Expression
	ThenBranch *BlockStatement
	ElseBranch *BlockStatement
}

func (expr *IfExpression) expressionNode()      {}
func (expr *IfExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(expr.Condition.String())
	out.WriteString(" ")
	out.WriteString(expr.ThenBranch.String())
	if expr.ElseBranch != nil {
		out.WriteString("else ")
		out.WriteString(expr.ElseBranch.String())
	}

	return out.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (stmt *BlockStatement) statementNode()       {}
func (stmt *BlockStatement) TokenLiteral() string { return stmt.Token.Literal }
func (stmt *BlockStatement) String() string {
	var out bytes.Buffer

	for i := range len(stmt.Statements) {
		out.WriteString(stmt.Statements[i].String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func (expr *FunctionLiteral) expressionNode()      {}
func (expr *FunctionLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *FunctionLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for i := range len(expr.Parameters) {
		params = append(params, expr.Parameters[i].String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(expr.Body.String())

	return out.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (expr *CallExpression) expressionNode()      {}
func (expr *CallExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for i := range len(expr.Arguments) {
		args = append(args, expr.Arguments[i].String())
	}

	out.WriteString(expr.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (expr *StringLiteral) expressionNode()      {}
func (expr *StringLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *StringLiteral) String() string       { return expr.Value }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (expr *ArrayLiteral) expressionNode()      {}
func (expr *ArrayLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *ArrayLiteral) String() string {
	var out bytes.Buffer

	var elements []string
	for i := range len(expr.Elements) {
		elements = append(elements, expr.Elements[i].String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (expr *IndexExpression) expressionNode()      {}
func (expr *IndexExpression) TokenLiteral() string { return expr.Token.Literal }
func (expr *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(expr.Left.String())
	out.WriteString("[")
	out.WriteString(expr.Index.String())
	out.WriteString("])")

	return out.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (expr *HashLiteral) expressionNode()      {}
func (expr *HashLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range expr.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type MacroLiteral struct {
	Token      token.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

func (expr *MacroLiteral) expressionNode()      {}
func (expr *MacroLiteral) TokenLiteral() string { return expr.Token.Literal }
func (expr *MacroLiteral) String() string {
	var out bytes.Buffer

	var params []string
	for i := range len(expr.Parameters) {
		params = append(params, expr.Parameters[i].String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(expr.Body.String())

	return out.String()
}
