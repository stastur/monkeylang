package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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
