package evaluator

import (
	"monkeylang/ast"
	"monkeylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.ReturnValue)}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.UnaryExpression:
		return evalUnaryExpression(node.Operator, Eval(node.Right))
	case *ast.BinaryExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalBinaryExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for i := range len(stmts) {
		result = Eval(stmts[i])

		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for i := range len(block.Statements) {
		result = Eval(block.Statements[i])

		if rv, ok := result.(*object.ReturnValue); ok {
			return rv
		}
	}

	return result
}

func nativeBoolToBooleanObject(b bool) *object.Boolean {
	if b {
		return TRUE
	}
	return FALSE
}

func evalUnaryExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusUnaryExpression(right)
	default:
		return nil
	}
}

func evalBangOperator(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusUnaryExpression(obj object.Object) object.Object {
	integer, ok := obj.(*object.Integer)
	if !ok {
		return NULL
	}

	return &object.Integer{Value: -integer.Value}
}

func evalBinaryExpression(op string, left, right object.Object) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		left := left.(*object.Integer)
		right := right.(*object.Integer)
		return evalIntegerBinaryExpression(op, left, right)
	}

	if op == "==" {
		return nativeBoolToBooleanObject(left == right)
	}

	if op == "!=" {
		return nativeBoolToBooleanObject(left != right)
	}

	return NULL
}

func evalIntegerBinaryExpression(op string, left, right *object.Integer) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch op {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return NULL
	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	if isTruthy(Eval(node.Condition)) {
		return Eval(node.ThenBranch)
	}

	if node.ElseBranch != nil {
		return Eval(node.ElseBranch)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case FALSE, NULL:
		return false
	default:
		return true
	}
}
