package main

import (
	"fmt"
	"log"
	"sort"
)

type ValueType int

const (
	Nothing ValueType = iota
	Bool
	Int
	Float
	Vector
	Tuple
	Function
	LastBuiltinType
)

func (v ValueType) String() string {
	if v < LastBuiltinType {
		return []string{"nothing", "bool", "int", "float", "vector", "tuple", "function", "nothing2"}[v]
	}
	return "User defined type"
}

type Scope struct {
	variables map[string]Value
}

func (result *Scope) Merge(other *Scope) {
	for k, v := range other.variables {
		//locally scoped varibales take precedence over those being merged in, that is if a key exists in result and other, the value in result will be used
		if _, in_result := result.variables[k]; !in_result {
			result.variables[k] = v
		}
	}
}

func (s *Scope) String() string {
	out := "{\n"

	keys := make([]string, 0)
	for key := range s.variables {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := s.variables[k]
		out += fmt.Sprintf("%s: %v\n", k, v)
	}
	return out + "}"
}

func EmptyScope() *Scope {
	return &Scope{
		variables: map[string]Value{},
	}
}

type ASTNode interface {
	Execute(r *Runtime)
	ReturnsType(r *Runtime) ValueType
}

var _ ASTNode = &SetNode{to: "", from: nil}
var _ ASTNode = &BoolLiteral{false}
var _ ASTNode = &IntLiteral{12}
var _ ASTNode = &PrintStatement{&IntLiteral{2}}
var _ ASTNode = &AddIntNode{}
var _ ASTNode = &TupleLiteral{}
var _ ASTNode = &AddAnyNode{}

type SetNode struct {
	to      string
	my_type ValueType
	from    ASTNode
}

func (sn *SetNode) Execute(r *Runtime) {
	sn.from.Execute(r)
	r.StackTop().variables[sn.to] = r.last_expression_result
}

func (sn *SetNode) ReturnsType(r *Runtime) ValueType {
	return sn.from.ReturnsType(r)
}

type GetNode struct {
	name   string
	v_type ValueType
}

func (gn *GetNode) Execute(r *Runtime) {
	r.last_expression_result = r.StackTop().variables[gn.name]

}

func (n *GetNode) ReturnsType(r *Runtime) ValueType {
	return n.v_type
}

type BoolLiteral struct {
	value bool
}

func (b *BoolLiteral) Execute(r *Runtime) {
	r.last_expression_result = &BoolType{
		name:  "",
		value: b.value,
	}
}
func (b *BoolLiteral) ReturnsType(r *Runtime) ValueType {
	return Bool
}

type IntLiteral struct {
	value int
}

func (*IntLiteral) ReturnsType(r *Runtime) ValueType {
	return Int
}

func (il *IntLiteral) Execute(r *Runtime) {

	r.last_expression_result = &IntType{
		name:  "",
		value: il.value,
	}
}

type AddAnyNode struct {
	left, right ASTNode
}

// Execute implements ASTNode
func (aan *AddAnyNode) Execute(r *Runtime) {
	operation, operation_exists := r.unary_operator_overloads[[2]ValueType{aan.left.ReturnsType(r), aan.right.ReturnsType(r)}]
	if !operation_exists {
		r.throwError(fmt.Sprintf("No overloaded operator exists between %v and %v", aan.left.ReturnsType(r), aan.right.ReturnsType(r)))
		return
	}
	aan.left.Execute(r)
	lval := r.last_expression_result
	aan.right.Execute(r)
	rval := r.last_expression_result
	r.last_expression_result = operation.operation(lval, rval)

}

// ReturnsType implements ASTNode
func (aan *AddAnyNode) ReturnsType(r *Runtime) ValueType {
	operation, operation_exists := r.unary_operator_overloads[[2]ValueType{aan.left.ReturnsType(r), aan.right.ReturnsType(r)}]
	if !operation_exists {
		// no corresponding operation defined
		return Nothing
	}
	return operation.ret_type
}

type AddIntNode struct {
	left, right ASTNode
}

func (ain *AddIntNode) Execute(r *Runtime) {
	ain.left.Execute(r)
	ln := r.last_expression_result
	ain.right.Execute(r)
	rn := r.last_expression_result

	l_int := 0
	r_int := 0
	switch li := ln.(type) {
	case *IntType:
		l_int = li.value
	}
	switch ri := rn.(type) {
	case *IntType:
		r_int = ri.value
	}
	sum := l_int + r_int
	r.last_expression_result = &IntType{
		name:  "",
		value: sum,
	}
}
func (ain *AddIntNode) ReturnsType(r *Runtime) ValueType {
	return Int
}

type SubIntNode struct {
	left, right ASTNode
}

func (sin *SubIntNode) Execute(r *Runtime) {
	sin.left.Execute(r)
	ln := r.last_expression_result
	sin.right.Execute(r)
	rn := r.last_expression_result

	l_int := 0
	r_int := 0
	switch li := ln.(type) {
	case *IntType:
		l_int = li.value
	}
	switch ri := rn.(type) {
	case *IntType:
		r_int = ri.value
	}
	sum := l_int + r_int
	r.last_expression_result = &IntType{
		name:  "",
		value: sum,
	}
}
func (sin *SubIntNode) ReturnsType(r *Runtime) ValueType {
	return Int
}

type TupleLiteral struct {
	values []ASTNode
}

func (tl *TupleLiteral) Execute(r *Runtime) {
	values := make([]Value, len(tl.values))
	for i := range tl.values {
		tl.values[i].Execute(r)
		values[i] = r.last_expression_result
	}
	r.last_expression_result = &TupleType{
		name:   "",
		values: values,
	}
}

func (*TupleLiteral) ReturnsType(r *Runtime) ValueType {
	return Tuple
}

type PrintStatement struct {
	argument ASTNode
}

func (ps *PrintStatement) Execute(r *Runtime) {
	ps.argument.Execute(r)
	switch arg := r.last_expression_result.(type) {
	case *IntType:
		fmt.Println(arg.value)
	case *TupleType:
		print_tuple(arg)
	default:
		log.Printf("Can not yet print type: %T: %v\n", arg, arg)
	}

}

func print_tuple(t *TupleType) {
	s := ""
	for i, v := range t.values {
		if v == nil {
			s += "nil"
		} else {
			s += v.String()
		}
		if i < len(t.values)-1 {
			s += " "
		}
	}
	s += ""
	fmt.Println(s)
}

func (*PrintStatement) ReturnsType(r *Runtime) ValueType {
	return Nothing
}

type FunctionDefinition struct {
	parameterNames []string
	parameterTypes []ValueType
	returnType     ValueType

	lines []ASTNode
}

func (fd *FunctionDefinition) Execute(r *Runtime) {
	r.NewIsolatedScope()

	r.PopScope()
}

type BinaryOperation struct {
	a_type, b_type ValueType
	ret_type       ValueType
	operation      func(a, b Value) Value
}
type Runtime struct {
	unary_operator_overloads map[[2]ValueType]BinaryOperation
	global_scope             *Scope
	scope_stack              []*Scope
	stack_depth              int
	last_expression_result   Value
	last_error               error

	ASTLines []ASTNode //outer level is []functions

	named_places map[string]int //index of function() in ASTLines

	current_line int
}

func (r *Runtime) StackTop() *Scope {
	return r.scope_stack[r.stack_depth]
}

/*
Creates a scope that has access to global variables, local variables outside it, and has space for new variables
var a int = 3

	for (var i int = 0; i<a; i++){
		//stuff in here has access to outer scope (a), global variables, and the new variables (i)
	}
*/
func (r *Runtime) NewLocalScope() {
	es := EmptyScope()
	last_stack := r.scope_stack[len(r.scope_stack)-1]
	es.Merge(last_stack) //last stack has global variables in it already

	r.scope_stack = append(r.scope_stack, es)
}

/*
Creates a scope that is isolated from everthing but global variables

	func stuff(args){
		in here only has args, global variables and anything else defined in here
	}
*/
func (r *Runtime) NewIsolatedScope() {
	es := EmptyScope()
	es.Merge(r.global_scope)
	r.scope_stack = append(r.scope_stack, es)
}
func (r *Runtime) PopScope() {
	r.scope_stack = r.scope_stack[:len(r.scope_stack)-1]
}

func (r *Runtime) throwError(s string) {
	log.Println(s)
}
func (r *Runtime) Run() {
	fmt.Println(r.ASTLines)
	for r.current_line < len(r.ASTLines) {
		fmt.Println("Line:", r.current_line)

		r.ASTLines[r.current_line].Execute(r)

		r.current_line++
		fmt.Println()
	}
}

func NewRuntime(program []ASTNode) *Runtime {
	return &Runtime{
		unary_operator_overloads: map[[2]ValueType]BinaryOperation{},
		global_scope:             &Scope{},
		scope_stack:              []*Scope{EmptyScope()},
		stack_depth:              0,
		last_expression_result:   nil,
		last_error:               nil,
		ASTLines:                 program,
		named_places:             map[string]int{},
		current_line:             0,
	}
}
