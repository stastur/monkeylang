package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
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
