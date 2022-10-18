package main

import (
	_ "embed"
	"fmt"
)

func main() {
	//	fmt.Println(src)

	/*
		program := []ASTNode{
			&SetNode{"a", Int, &IntLiteral{value: 12}},                                  //a=12
			&SetNode{"a", Int, &AddIntNode{&GetNode{"a", Int}, &IntLiteral{value: 11}}}, //a = a + 11
			&PrintStatement{&GetNode{"a", Int}},                                         //print a
		}
	*/
	/*
		program2 := []ASTNode{
			&SetNode{"a", Int, &IntLiteral{value: 12}},      //var a int = 12
			&SetNode{"b", Int, &IntLiteral{value: 6}},       //var a int = 6
			&SetNode{"c", Bool, &BoolLiteral{value: false}}, // var c bool = false
			&SetNode{"d", Tuple, &TupleLiteral{[]ASTNode{ // var d = (a b c)
				&GetNode{"a", Int},
				&GetNode{"b", Int},
				&GetNode{"c", Bool},
			}}},
			&PrintStatement{&GetNode{"d", Tuple}}, // print d
		}
		runtime := NewRuntime(program2)
		runtime.Run()
	*/
	//line_src := "var a int = 12"
	//line_src := "var a string = \"w\""
	//line_src := "var a vec<int> = [1 2 3 4]"
	//line_src := "var a var = 0"
	line_src := "var abcde = 0"
	fmt.Println(line_src)
	toks := Tokenize(line_src)
	MakeTree(toks, line_src)
	fmt.Println(toks)

}
