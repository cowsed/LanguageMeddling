package main

import (
	"fmt"
	"sort"
	"strings"
)

type ErrorCollector struct {
	errs       []error
	shouldstop bool
}

func (ec *ErrorCollector) SayErrors(lines []string) {
	for _, err := range ec.errs {
		switch e := err.(type) {
		case LocatedError:
			e.line_src = lines[e.line-1]
			err = e
		}
		fmt.Println(err.Error())

	}
}
func (ec *ErrorCollector) ShouldStop() {
	ec.shouldstop = true
}
func (ec *ErrorCollector) AddError(err error) {
	ec.errs = append(ec.errs, err)
}

func NewLocatedError(line, index int, msg string) LocatedError {
	return LocatedError{
		line:  line,
		index: index,
		msg:   msg,
	}
}

type LocatedError struct {
	line  int
	index int
	msg   string

	line_src string
}

func (le LocatedError) Error() string {
	s := ""
	if le.line_src != "" {
		s += le.line_src + "\n"
		s += strings.Repeat(" ", le.index) + "^\n"
		s += fmt.Sprintf("%s", le.msg)

	} else {
		s += fmt.Sprintf("line %d:%d : %s", le.line, le.index, le.msg)
	}
	return s
}

// allows parsers to tell the system that there may be a problem here in the future
// example var x structA could be true or false depending on whether or not structA is defined in the future
type ParseChecker struct {
	ErrorCollector
	src_lines         []string
	num_defined_types int

	type_nums            map[string]int
	types_defined        map[string]bool
	declared_type_checks map[string][]TypeDefinedCheck
}

func (pc *ParseChecker) GetTypeNum(type_name string) int {
	if num, exists := pc.type_nums[type_name]; exists {
		//seen before, already have a number for it
		return num
	} else {
		//never before seen, need to create a number for it
		return pc.AddType(type_name)
	}
}
func (pc *ParseChecker) AddType(name string) int {
	pc.num_defined_types++
	type_num := pc.num_defined_types + int(LastBuiltinType)
	pc.type_nums[name] = type_num
	return type_num
}

/*
func (pc *ParseChecker) DefineType(type_name string, at_line int) {
	if _, type_exists := pc.type_nums[type_name]; type_exists {
		//redifinition of type type_name
		pc.AddError(NewLocatedError(at_line, 0, "type redeclaration"))
	} else {
		// add to known types
		pc.num_defined_types++
		type_num := pc.num_defined_types + int(LastBuiltinType)
		pc.type_nums[type_name] = type_num
		//remove all the type listeners
		delete(pc.declared_type_checks, type_name)
	}
}
*/
// add a watcher that will throw an error if the type is not defined by the end of analysis
func (pc *ParseChecker) EnsureTypeDefined(tdc TypeDefinedCheck) {
	//if type already defined, dont add watcher
	type_name := tdc.type_name
	if already_defined := pc.types_defined[type_name]; already_defined {
		fmt.Println("type", type_name, "already defined")
		return
	}
	// not yet defined, add watcher - 2 options if its already in the map
	others, already_in := pc.declared_type_checks[type_name]
	if already_in {
		fmt.Println("type", type_name, "already has a watcher")

		others = append(others, tdc)
		pc.declared_type_checks[type_name] = others
	} else {
		fmt.Println("type", type_name, "getting set")

		pc.declared_type_checks[type_name] = []TypeDefinedCheck{tdc}
	}

}

func (pc *ParseChecker) SayErrors() {
	//if there are any errors not resolved in the analysis section, speak now or forever hold your peace
	fmt.Println("\n\nAnalysis Errors:")
	undefined_types_keys := make([]string, len(pc.declared_type_checks))
	fmt.Println("declared type checks:", pc.declared_type_checks)
	i := 0
	for k, _ := range pc.declared_type_checks {
		undefined_types_keys[i] = k
		i++
	}
	fmt.Println(undefined_types_keys)
	sort.Strings(undefined_types_keys)

	for _, key := range undefined_types_keys {
		for _, e := range pc.declared_type_checks[key] {
			e.error_if_not.line_src = pc.src_lines[e.error_if_not.line-1]
			fmt.Println(e.error_if_not.Error())
		}
	}
	pc.ErrorCollector.SayErrors(pc.src_lines)
}

type TypeDefinedCheck struct {
	type_name    string
	error_if_not LocatedError
}

func MakeTree(token_lines [][]Token, src string) {
	lines := strings.Split(src, "\n")
	pc := &ParseChecker{
		src_lines:            lines,
		ErrorCollector:       ErrorCollector{},
		num_defined_types:    0,
		type_nums:            map[string]int{},
		declared_type_checks: map[string][]TypeDefinedCheck{},
	}
	var ast_head []ASTNode = []ASTNode{}

	for _, line := range token_lines {
		tg := &TokenGiver{toks: line, index: 0}
		tok := tg.PeekNext()
		switch tok.TokenType {
		case Var_TType:
			var_ast := TreeifyVarStatement(tg, pc)
			ast_head = append(ast_head, var_ast...)
		}

	}
	for i := range ast_head {
		fmt.Printf("%+v\n", ast_head[i])
	}
	pc.SayErrors()

}

// returns true, intenal if it is vec<internal>, else false ""
func is_vector_wrapper(s string) (bool, string) {
	prefix := "vector<"
	if s[0:len(prefix)] == prefix && s[len(s)-2:] == ">" {
		internal := s[len(prefix) : len(s)-1]
		return true, internal
	}
	return false, ""
}

func TreeifyVarStatement(tg *TokenGiver, pc *ParseChecker) []ASTNode {
	nodes := []ASTNode{}
	var_tok := tg.ConsumeNext() // should just be var
	if !tg.HasNext() {
		pc.AddError(NewLocatedError(var_tok.line, var_tok.index_end, "expected variable name"))
		return nodes
	}
	name_tok := tg.ConsumeNext()
	if !tg.HasNext() || (tg.PeekNext().TokenType != BuiltinType_TType && tg.PeekNext().TokenType != Name_TType) {
		pc.AddError(NewLocatedError(var_tok.line, name_tok.index_end, "expected variable type"))
		return nodes
	}
	var_type_tok := tg.ConsumeNext()
	var actual_type ValueType
	var is_primitive = true
	var is_simple = true //not a vec or tuple
	if var_type_tok.TokenType == BuiltinType_TType {
		switch var_type_tok.text {
		case "bool":
			actual_type = Bool
		case "int":
			actual_type = Int
		case "float":
			actual_type = Float
		case "string":
			actual_type = String
		default:
			//filter out complex types
			if is_vec, sub_type := is_vector_wrapper(var_type_tok.text); is_vec {
				actual_type = Vector
				is_primitive = false
				fmt.Println(sub_type)
				panic("unimplemented")
			} else {
				pc.AddError(NewLocatedError(var_tok.line, var_type_tok.index_start, "unknown builtin type, this should probably never happen if this analysis is well written"))
				actual_type = NoType
			}
		}
	}
	if var_type_tok.TokenType == Name_TType { //user defined type
		is_primitive = false
		is_simple = true
		type_name := var_type_tok.text

		type_num := pc.GetTypeNum(type_name)
		//add watcher to make sure this type actually gets defined later
		pc.EnsureTypeDefined(TypeDefinedCheck{
			type_name: var_type_tok.text,
			error_if_not: LocatedError{
				line:  var_tok.line,
				index: var_type_tok.index_end,
				msg:   fmt.Sprintf("type %s was never defined", type_name),
			},
		})
		actual_type = ValueType(type_num)
	}
	if is_primitive {
		nodes = append(nodes, &DeclareNode{
			name:    name_tok.text,
			my_type: actual_type,
		})
	} else if is_simple { //not array type
		nodes = append(nodes, &DeclareNode{
			name:    name_tok.text,
			my_type: actual_type,
		})
	} else {
		panic("unimplemented")
	}

	if !tg.HasNext() {
		//we good, just a declaration, not a setting
		return nodes
	}
	if tg.PeekNext().TokenType != Assignment {
		pc.AddError(NewLocatedError(var_tok.line, var_type_tok.index_start, "expected `=` or newline"))
		pc.ShouldStop()
		return nodes
	}
	tg.ConsumeNext() //take =

	exp := TreeifyExpression(tg, pc)

	nodes = append(nodes, &SetNode{
		to:      name_tok.text,
		my_type: actual_type,
		from:    exp,
	})
	//panic("unimplemented declaration and assignment in the same line")

	return nodes
}
func TreeifyExpression(tg *TokenGiver, pc *ParseChecker) ASTNode {
	return &IntLiteral{
		value: -1,
	}
}

type TokenGiver struct {
	toks  []Token
	index int
}

func (tg *TokenGiver) HasNext() bool {
	return tg.index < len(tg.toks)
}

func (tg *TokenGiver) PeekNext() Token {
	return tg.toks[tg.index]
}
func (tg *TokenGiver) ConsumeNext() Token {
	tg.index++
	return tg.toks[tg.index-1]

}
