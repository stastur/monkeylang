package evaluator

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: Eval(node.ReturnValue)}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.UnaryExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(node.Operator, right)
	case *ast.BinaryExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

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

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for i := range len(block.Statements) {
		result = Eval(block.Statements[i])

		switch result := result.(type) {
		case *object.ReturnValue, *object.Error:
			return result
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
		return newError("unknown operator: %s%s", op, right.Type())
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
		return newError("unknown operator: -%s", obj.Type())
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

	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s",
			left.Type(), op, right.Type())
	}

	return newError("unknown operator: %s %s %s",
		left.Type(), op, right.Type())
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
		return newError("unknown operator: %s %s %s",
			op, left.Type(), right.Type())
	}
}

func evalIfExpression(node *ast.IfExpression) object.Object {
	cond := Eval(node.Condition)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
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
