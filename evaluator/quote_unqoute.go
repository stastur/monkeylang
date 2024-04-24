package evaluator

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/object"
	"monkeylang/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		expr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if expr.Function.TokenLiteral() != "unquote" {
			return node
		}

		if len(expr.Arguments) > 1 {
			return node
		}

		return convertObjectToAstNode(Eval(expr.Arguments[0], env))
	})
}

func isUnquoteCall(node ast.Node) bool {
	expr, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}

	return expr.Function.TokenLiteral() == "unquote"
}

func convertObjectToAstNode(obj object.Object) ast.Node {
	// TODO: Handle other object types
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	case *object.Quote:
		return obj.Node

	default:
		return nil
	}
}
