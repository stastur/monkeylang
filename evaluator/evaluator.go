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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.IdentifierExpression:
		return evalIdentifier(node, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.UnaryExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(node.Operator, right)
	case *ast.BinaryExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalBinaryExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.FunctionLiteral:
		return &object.Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args, err := evalExpressions(node.Arguments, env)
		if err != nil {
			return err
		}

		return applyFunction(function, args)
	}

	return nil
}

func evalExpressions(
	exprs []ast.Expression,
	env *object.Environment,
) ([]object.Object, *object.Error) {
	var result []object.Object

	for i := range len(exprs) {
		evaluated := Eval(exprs[i], env)
		if isError(evaluated) {
			return nil, evaluated.(*object.Error)
		}
		result = append(result, evaluated)
	}

	return result, nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for i := range len(stmts) {
		result = Eval(stmts[i], env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for i := range len(block.Statements) {
		result = Eval(block.Statements[i], env)

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

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(node.ThenBranch, env)
	}

	if node.ElseBranch != nil {
		return Eval(node.ElseBranch, env)
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

func evalIdentifier(
	node *ast.IdentifierExpression,
	env *object.Environment,
) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

func applyFunction(
	obj object.Object,
	args []object.Object,
) object.Object {
	function, ok := obj.(*object.Function)
	if !ok {
		return newError("not a function: %s", obj.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)

	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(
	fn *object.Function,
	args []object.Object,
) *object.Environment {
	env := object.NewEnclosedEnv(fn.Env)

	for i := range len(fn.Parameters) {
		if i < len(args) {
			env.Set(fn.Parameters[i].Value, args[i])
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rv, ok := obj.(*object.ReturnValue); ok {
		return rv.Value
	}

	return obj
}
