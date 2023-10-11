package ast

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode"
	"unicode/utf8"
)

func lexFile(path string) tokens {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !utf8.Valid(data) {
		fmt.Println(path, " is not valid utf8\n")
		os.Exit(1)
	}
	return lexBytes(data)
}

func lexBytes(b []byte) tokens {
	s := &subject{input: b, line_no: 1}
	s.lex() // puts the first token into s.current for initialization
	return s
}

func ShowLex(path string) {
	ts := lexFile(path)
	for ts.peek().typeOf != eof {
		fmt.Println(ts.take())
	}
}

// the ast will be built via these two primitives peek and take
// concrete type &subject implements this interface
// eof token means end of stream, and repeated take()s will return eof forever
type tokens interface {
	peek() token
	take() token
}

// token type definitions
// since token is the return value for peek and take, it also serves
// as part of the interface for the parser by inspecting the kind

type tokType int

const (
	iden tokType = iota // iden must be zero; keywords["foo"]==0 implies iden
	lambda
	dot
	lpar
	rpar
	define
	let
	symbol
	lbrack
	rbrack
	match
	pipe
	backslash
	wildcard
	number
	ifexpr
	booltrue
	boolfalse
	eof
)

// make printable
var showKind = [...]string{
	"iden",
	"lambda",
	"dot",
	"lpar",
	"rpar",
	"define",
	"let",
	"symbol",
	"lbrack",
	"rbrack",
	"match",
	"pipe",
	"backslash",
	"wildcard",
	"number",
	"ifexpr",
	"booltrue",
	"boolfalse",
	"eof",
}

// these are "seperators", they can be right next to each other in the source
var sepChars = map[string]tokType{
	"λ": lambda,
	".": dot,
	"(": lpar,
	")": rpar,
	"[": lbrack,
	"]": rbrack,
	"\\": backslash, // for pattern variables rn
}

var key_words = map[string]tokType{
	"def":   define,
	"let":   let,
	"?":     ifexpr,
	"true":  booltrue,
	"false": boolfalse,
	"μ":     match,
	"|":     pipe,
	"_":     wildcard,
}

type token struct {
	typeOf  tokType
	literal string
	line_no int
}

// to satisfy the interface to the parser
func (s *subject) peek() token {
	return s.current
}
func (s *subject) take() token {
	c := s.current
	s.lex()
	return c
}

// from here down, code is internal to the lexing phase
type subject struct {
	input   []byte
	begin   int
	end     int
	current token
	line_no int
}

// core process for the scanning of the text
func (s *subject) lex() {
	r := s.peek_rune()
	c := string(r)
	s.step()

	switch {
	case r == utf8.RuneError:
		s.commit(eof)

	// white space
	case c == "\n":
		s.line_no++
		s.junk()
	case unicode.IsSpace(r):
		s.junk()
	case c == ";":
		s.lexes(func(r rune) bool { return !(string(r) == "\n") })
		s.junk()

	// seps
	case sepChars[c] != 0: // this works, since idens are the zero value
		s.commit(sepChars[c])

	// symbolic constants
	case c == ":":
		s.lexes(is_char)
		s.commit(symbol)
	
		
	// words
	case is_char(r):
		s.lexes(is_char)
		word := s.inspect()
		if isNumber(word) {
			s.commit(number)
		} else {
			s.commit(key_words[s.inspect()]) // commits iden when keyw==0
		}
	}
}

func (s *subject) lexes(f func(rune) bool) {
	for s.peek_rune() != utf8.RuneError && f(s.peek_rune()) {
		s.step()
	}
}


func isNumber(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// primitive scanning functions
func (s *subject) peek_rune() rune {
	r, _ := utf8.DecodeRune(s.input[s.end:])
	return r
}
func (s *subject) step() {
	_, w := utf8.DecodeRune(s.input[s.end:])
	s.end += w
}
func (s *subject) junk() {
	s.begin = s.end
	s.lex()
}
func (s *subject) inspect() string {
	return string(s.input[s.begin:s.end])
}
func (s *subject) commit(tk tokType) {
	s.current = token{tk, s.inspect(), s.line_no}
	s.begin = s.end
}

// isgraphic ensures less strange unicode stuff for chars
func is_char(r rune) bool {
	_, is_sep := sepChars[string(r)]
	return unicode.IsGraphic(r) && !is_sep && !unicode.IsSpace(r)
}

// show tokens for inspection

func (t tokType) String() string {
	return fmt.Sprintf(showKind[t])
}
func (t token) String() string {
	msg := "line %v | %v | %v "
	return fmt.Sprintf(msg, t.line_no, showKind[t.typeOf], t.literal)
}
