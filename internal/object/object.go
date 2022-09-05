package object

import "fmt"

type ObjectType string

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (o *Integer) Inspect() string {
	return fmt.Sprintf("%d", o.Value)
}
func (o *Integer) Type() ObjectType {
	return IntegerObj
}

type Boolean struct {
	Value bool
}

func (o *Boolean) Inspect() string {
	return fmt.Sprintf("%t", o.Value)
}
func (o *Boolean) Type() ObjectType {
	return BooleanObj
}

type Null struct{}

func (o *Null) Inspect() string {
	return "null"
}
func (o *Null) Type() ObjectType {
	return NullObj
}

type ReturnValue struct {
	Value Object
}

func (o *ReturnValue) Type() ObjectType {
	return ReturnValueObj
}
func (o *ReturnValue) Inspect() string {
	return o.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ErrorObj
}
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}
