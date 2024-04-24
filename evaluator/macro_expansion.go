package evaluator

import (
	"monkeylang/ast"
	"monkeylang/object"
)

func DefineMacros(program *ast.Program, env *object.Environment) {
	statements := make([]ast.Statement, 0, len(program.Statements))

	for i := range program.Statements {
		stmt := program.Statements[i]

		if isMacroDefinition(stmt) {
			addMacro(stmt, env)
		} else {
			statements = append(statements, stmt)
		}
	}

	program.Statements = statements
}

func isMacroDefinition(stmt ast.Statement) bool {
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		return false
	}

	_, ok = letStmt.Value.(*ast.MacroLiteral)

	return ok
}

func addMacro(stmt ast.Statement, env *object.Environment) {
	letStmt := stmt.(*ast.LetStatement)
	macroLiteral := letStmt.Value.(*ast.MacroLiteral)

	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Body:       macroLiteral.Body,
		Env:        env,
	}

	env.Set(letStmt.Name.Value, macro)
}

func ExpandMacros(program *ast.Program, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpr, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		macro, ok := isMacroCall(callExpr, env)
		if !ok {
			return node
		}

		args := quoteArgs(callExpr)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)

		quote, ok := evaluated.(*object.Quote)
		if !ok {

			panic("we only support returning AST-nodes from macros")
		}

		return quote.Node
	})
}

func isMacroCall(
	expr *ast.CallExpression,
	env *object.Environment,
) (*object.Macro, bool) {
	ident, ok := expr.Function.(*ast.IdentifierExpression)
	if !ok {
		return nil, false
	}

	obj, ok := env.Get(ident.Value)
	if !ok {
		return nil, false
	}

	macro, ok := obj.(*object.Macro)
	return macro, ok
}

func quoteArgs(expr *ast.CallExpression) []*object.Quote {
	args := make([]*object.Quote, len(expr.Arguments))

	for i := range expr.Arguments {
		args[i] = &object.Quote{Node: expr.Arguments[i]}
	}

	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote,
) *object.Environment {
	extended := object.NewEnclosedEnv(macro.Env)

	for paramIdx, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIdx])
	}
	return extended
}
