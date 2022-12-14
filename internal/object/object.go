package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/vancanhuit/monkey/internal/ast"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

const (
	IntegerObj     = "INTEGER"
	BooleanObj     = "BOOLEAN"
	NullObj        = "NULL"
	ReturnValueObj = "RETURN_VALUE"
	ErrorObj       = "ERROR"
	FunctionObj    = "FUNCTION"
	StringObj      = "STRING"
	BuiltinObj     = "BUILTIN"
	ArrayObj       = "ARRAY"
	HashObj        = "HASH"
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

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (o *Function) Type() ObjectType {
	return FunctionObj
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

type String struct {
	Value string
}

func (o *String) Type() ObjectType {
	return StringObj
}
func (o *String) Inspect() string {
	return o.Value
}

type Builtin struct {
	Fn BuiltinFunction
}

func (o *Builtin) Type() ObjectType {
	return BuiltinObj
}
func (o *Builtin) Inspect() string {
	return "builtin function"
}

type Array struct {
	Elements []Object
}

func (o *Array) Type() ObjectType {
	return ArrayObj
}
func (o *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range o.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HashObj
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(
			pairs,
			fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
