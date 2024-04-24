package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkeylang/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	QUOTE_OBJ        = "QUOTE"
	UNQUOTE_OBJ      = "UNQUOTE"
	MACRO_OBJ        = "MACRO"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (o *Integer) Inspect() string  { return fmt.Sprintf("%d", o.Value) }
func (o *Integer) Type() ObjectType { return INTEGER_OBJ }

type Boolean struct {
	Value bool
}

func (o *Boolean) Inspect() string  { return fmt.Sprintf("%t", o.Value) }
func (o *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Null struct{}

func (o *Null) Inspect() string  { return "null" }
func (o *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Inspect() string  { return o.Value.Inspect() }
func (o *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
}

func (o *Error) Inspect() string  { return "Error: " + o.Message }
func (o *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameters []*ast.IdentifierExpression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (o *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range o.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(o.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (o *Function) Type() ObjectType { return FUNCTION_OBJ }

type String struct {
	Value string
}

func (o *String) Inspect() string  { return o.Value }
func (o *String) Type() ObjectType { return STRING_OBJ }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (o *Builtin) Inspect() string  { return "BUILTIN" }
func (o *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Array struct {
	Elements []Object
}

func (o *Array) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[\n")
	for i := range len(o.Elements) {
		out.WriteString(o.Elements[i].Inspect())
		out.WriteString(",\n")
	}
	out.WriteString("]")

	return out.String()
}
func (o *Array) Type() ObjectType { return ARRAY_OBJ }

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (o *Boolean) HashKey() HashKey {
	var value uint64

	if o.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: o.Type(), Value: value}
}

func (o *Integer) HashKey() HashKey {
	return HashKey{Type: o.Type(), Value: uint64(o.Value)}
}

func (o *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(o.Value))

	return HashKey{Type: o.Type(), Value: uint64(h.Sum64())}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (o *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs,
			fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type Quote struct {
	Node ast.Node
}

func (o *Quote) Type() ObjectType { return QUOTE_OBJ }
func (o *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", o.Node.String())
}

type Macro struct {
	Parameters []*ast.IdentifierExpression
	Body       *ast.BlockStatement
	Env        *Environment
}

func (o *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range o.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(o.Body.String())
	out.WriteString("\n}")
	return out.String()
}
func (o *Macro) Type() ObjectType { return MACRO_OBJ }
