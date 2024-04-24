package ast

type ModifierFunc func(Node) Node

func Modify(node Node, modifier ModifierFunc) Node {
	switch node := node.(type) {
	// TODO: error handling when type casting
	case *Program:
		for i := range node.Statements {
			node.Statements[i] = Modify(node.Statements[i], modifier).(Statement)
		}

	case *ExpressionStatement:
		node.Expression = Modify(node.Expression, modifier).(Expression)

	case *BinaryExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Right = Modify(node.Right, modifier).(Expression)

	case *UnaryExpression:
		node.Right = Modify(node.Right, modifier).(Expression)

	case *IndexExpression:
		node.Left = Modify(node.Left, modifier).(Expression)
		node.Index = Modify(node.Index, modifier).(Expression)

	case *IfExpression:
		node.Condition = Modify(node.Condition, modifier).(Expression)
		node.ThenBranch = Modify(node.ThenBranch, modifier).(*BlockStatement)
		if node.ElseBranch != nil {
			node.ElseBranch = Modify(node.ElseBranch, modifier).(*BlockStatement)
		}

	case *BlockStatement:
		for i := range node.Statements {
			node.Statements[i] = Modify(node.Statements[i], modifier).(Statement)
		}

	case *CallExpression:
		node.Function = Modify(node.Function, modifier).(Expression)
		for i, expr := range node.Arguments {
			node.Arguments[i] = Modify(expr, modifier).(Expression)
		}

	case *ReturnStatement:
		node.ReturnValue = Modify(node.ReturnValue, modifier).(Expression)

	case *LetStatement:
		node.Value = Modify(node.Value, modifier).(Expression)

	case *FunctionLiteral:
		for i := range node.Parameters {
			node.Parameters[i] =
				Modify(node.Parameters[i], modifier).(*IdentifierExpression)
		}
		node.Body = Modify(node.Body, modifier).(*BlockStatement)

	case *ArrayLiteral:
		for i := range node.Elements {
			node.Elements[i] = Modify(node.Elements[i], modifier).(Expression)
		}

	case *HashLiteral:
		newPairs := make(map[Expression]Expression, len(node.Pairs))
		for k, v := range node.Pairs {
			newKey := Modify(k, modifier).(Expression)
			newValue := Modify(v, modifier).(Expression)
			newPairs[newKey] = newValue
		}
	}

	return modifier(node)
}
