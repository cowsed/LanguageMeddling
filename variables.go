package main

import "fmt"

type Value interface {
	Name() string
	String() string
	Type() ValueType
}

var _ Value = &BoolType{name: "a", value: false}
var _ Value = &IntType{name: "a", value: 20}
var _ Value = &TupleType{}

type BoolType struct {
	name  string
	value bool
}

func (b *BoolType) Name() string {
	return b.name
}
func (b *BoolType) String() string {
	return fmt.Sprint(b.value)
}

func (b *BoolType) Type() ValueType {
	return Bool
}

// an int variable
type IntType struct {
	name  string
	value int
}

func (*IntType) Type() ValueType {
	return Int
}
func (i *IntType) Name() string {
	return i.name
}
func (i *IntType) String() string {
	return fmt.Sprint(i.value)
}

type TupleType struct {
	name   string
	values []Value
}

func (tt *TupleType) Name() string {

	return tt.name
}
func (tt *TupleType) String() string {
	return fmt.Sprint(tt.values)
}

// Type implements Value
func (*TupleType) Type() ValueType {
	return Tuple
}
