package main

import (
	"fmt"
	"log"
	"strings"
)

func Tokenize(src_txt string) [][]Token {
	lines := strings.Split(src_txt, "\n")
	token_lines := make([][]Token, 0, len(lines)/2) //safe bet that at least half of all lines are code not whitespace, capacity not length tho
	for i := range lines {
		lt := LineTokenizer{
			line_src: lines[i],
			index:    0,
			line_num: i + 1,
		}
		toks := lt.Parse()
		token_lines = append(token_lines, toks)
	}
	return token_lines
}

type LineTokenizer struct {
	line_src string
	index    int
	line_num int
}

func (lp *LineTokenizer) throwError(msg string, index int, stop_line bool) {
	log.Println(msg, " at ", "will stop? ", stop_line)
}
func (lp *LineTokenizer) Rest() string {
	t := lp.line_src[lp.index : len(lp.line_src)-1]
	lp.index = len(lp.line_src)
	return t
}
func (lp *LineTokenizer) ParseNumber(initial string) string {
	sofar := initial
	for lp.HasNext() {
		next := lp.PeekNext()
		if !strings.Contains("1234567890e.", next) {
			break
		} else {
			sofar += lp.ConsumeNext()
		}
	}
	return sofar
}
func (lp *LineTokenizer) ParseText(initial string) string {
	sofar := initial
	for lp.HasNext() {
		next := lp.PeekNext()
		if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_", next) {
			break
		} else {
			sofar += lp.ConsumeNext()
		}
	}
	return sofar
}
func (lp *LineTokenizer) ParseQuotedText() string {
	sofar := ""
	for lp.HasNext() {
		next := lp.PeekNext()
		if next == "\"" {
			lp.ConsumeNext()
			return sofar
		} else if next == "\\" {
			lp.ConsumeNext() // \
			special := lp.ConsumeNext()
			switch special {
			case "n":
				sofar += "\n"
			case "t":
				sofar += "\t"
			}
		} else {
			sofar += lp.ConsumeNext()
		}
	}
	//ran out of text and no ""
	lp.throwError("no closing \"", lp.index, true)
	return ""
}

func (lp *LineTokenizer) PeekNext() string {
	return lp.line_src[lp.index : lp.index+1]
}
func (lp *LineTokenizer) ConsumeNext() string {
	s := lp.line_src[lp.index : lp.index+1]
	lp.index++
	return s
}
func (lp *LineTokenizer) HasNext() bool {
	return lp.index < len(lp.line_src)
}

func (lp *LineTokenizer) Parse() []Token {
	toks := []Token{}
	for lp.HasNext() {
		start := lp.index
		tok := Token{}
		s := lp.ConsumeNext()
		switch s {
		case "!": //only ever appears in a one long token
			tok = Token{TokenType: 0, text: "", index_start: start, index_end: start + 1}
		case "=": //= ==
			next := lp.PeekNext()
			if next == "=" {
				//is ==
				lp.ConsumeNext()
				tok = Token{TokenType: Equality, text: "==", index_start: start, index_end: start + 2}
			} else {
				//just =
				tok = Token{
					TokenType: Assignment, text: "=", index_start: start, index_end: start + 1,
				}
			}
		case "&":
			next := lp.PeekNext()
			if next == "&" { //&& and
				tok = Token{TokenType: And, text: "&&", index_start: start, index_end: start + 2}
			} else { //& reference of
				tok = Token{TokenType: Reference, text: "&", index_start: start, index_end: start + 1}
			}
		case "|":
			//no such thing as |, only || else error
			next := lp.PeekNext()
			if next != "|" {
				lp.throwError("no such operator `|`; did you mean `||`?", start, false)
				continue
			} else {
				lp.ConsumeNext()
				tok = Token{TokenType: Or, text: "||", index_start: start, index_end: start + 2}
			}
		case "+":
			tok = Token{TokenType: Plus, text: "+", index_start: start, index_end: start + 1}
		case "*":
			tok = Token{TokenType: Multiply, text: "*", index_start: start, index_end: start + 1}
		case "/":
			next := lp.PeekNext()
			if next == "/" { //is a comment
				lp.ConsumeNext() //get rid of next /
				comment_src := lp.Rest()
				tok = Token{TokenType: Comment_TType, text: comment_src, index_start: start, index_end: lp.index}
			} else { //is division
				tok = Token{TokenType: Divide, text: "/", index_start: start, index_end: start + 1}
			}
		case ",":
			tok = Token{TokenType: Comma, text: ",", index_start: start, index_end: start + 1}
		case ".":
			next := lp.PeekNext()
			if strings.Contains("1234567890", next) {
				//is a num literal starting with . ie. .2
				src := lp.ParseNumber(next)
				tok = Token{TokenType: NumLiteral_TType, text: src, index_start: start, index_end: lp.index}
			} else {
				//is a dot
				tok = Token{TokenType: Dot, text: ".", index_start: start, index_end: start + 1}
			}
		case "(":
			tok = Token{TokenType: OpenParen, text: "(", index_start: start, index_end: start + 1}
		case ")":
			tok = Token{TokenType: CloseParen, text: ")", index_start: start, index_end: start + 1}

		case "<":
			tok = Token{TokenType: OpenAlligator, text: "<", index_start: start, index_end: start + 1}
		case ">":
			tok = Token{TokenType: CloseAlligator, text: ">", index_start: start, index_end: start + 1}

		case "[":
			tok = Token{TokenType: OpenSquare, text: "[", index_start: start, index_end: start + 1}
		case "]":
			tok = Token{TokenType: CloseParen, text: "]", index_start: start, index_end: start + 1}

		case "{":
			tok = Token{TokenType: OpenCurly, text: "{", index_start: start, index_end: start + 1}
		case "}":
			tok = Token{TokenType: CloseCurly, text: "}", index_start: start, index_end: start + 1}

		case "\"":
			txt := lp.ParseQuotedText()
			tok = Token{TokenType: StringLiteral_TType, text: txt, index_start: start, index_end: start + lp.index}
		case "-":

			//is subtraction - if its a negative numbere, that will be taken care of when making the tree(need knowledge about the last token , if it was an operator then we take this to be negative, if its standalone its negate, and if its after an operator its actually minus)0
			tok = Token{TokenType: Minus, text: "-", index_start: start, index_end: start + 1}

		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			src := lp.ParseNumber(s)
			tok = Token{TokenType: NumLiteral_TType, text: src, index_start: start, index_end: lp.index}
		case "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z":
			//is the start of text
			txt := lp.ParseText(s)
			tok = TextToken(txt)
			tok.index_start = start
			tok.index_end = lp.index
		case " ", "\t":
			continue
		default:
			fmt.Println("unknown char", s)
			lp.throwError("Unknown character", lp.index, true)
		}

		tok.line = lp.line_num
		toks = append(toks, tok)
	}
	return toks
}
func TextToken(txt string) Token {
	switch txt {
	case "var":
		return Token{TokenType: Var_TType, text: txt}
	case "print":
		return Token{TokenType: Print_TType, text: txt}
	case "int":
		return Token{TokenType: BuiltinType_TType, text: txt}
	case "float":
		return Token{TokenType: BuiltinType_TType, text: txt}
	case "string":
		return Token{TokenType: StringLiteral_TType, text: txt}
	case "vec":
		return Token{TokenType: Vec_TType, text: txt}
	}
	return Token{
		TokenType: Name_TType,
		text:      txt,
	}
}

type Token struct {
	TokenType
	text                   string
	line                   int
	index_start, index_end int
}

func (t Token) String() string {
	return fmt.Sprintf("%s:%s", &t.TokenType, t.text)
}
func (t TokenType) String() string {
	names := []string{"Unknown_TType", "Var_TType", "Name_TType", "NumLiteral_TType", "StringLiteral_TType", "Vec_TType", "IntType_TType", "FloatType_TType", "StringType_TType", "Print_TType", "Comment_TType", "OpenAlligator", "CloseAlligator", "OpenParen", "CloseParen", "OpenCurly", "CloseCurly", "OpenSquare", "CloseSquare", "Assignment", "Equality", "Plus", "Minus", "Multiply", "Divide", "Reference", "Not", "Or", "And"}
	return names[t]
}

/*
var a int = 12
Var Name IntType Assignment IntLiteral
*/

type TokenType int

const (
	Unknown_TType       TokenType = iota
	Var_TType                     //var
	Name_TType                    // var_name
	NumLiteral_TType              // 1, 2, -4 , 1e23, 0.231
	StringLiteral_TType           //"wow"
	Vec_TType                     //vec
	BuiltinType_TType             //int, string, etc
	Print_TType                   //print
	Comment_TType                 // //
	//Brackets
	OpenAlligator
	CloseAlligator
	OpenParen
	CloseParen
	OpenCurly
	CloseCurly
	OpenSquare
	CloseSquare
	Comma
	Dot
	//Operators
	Assignment //=
	Equality   //==
	Plus       //+
	Minus      //-
	Multiply   //*
	Divide     // /
	Reference  //&
	Not        //!
	Or         //||
	And        //&&

)
