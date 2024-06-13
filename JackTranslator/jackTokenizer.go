package JackTranslator

import (
	"bufio"
	"os"
)

// throws if illegal token
// use lexcial elements section from slides

// Constants for token types
const (
	KEYWORD      = 1
	SYMBOL       = 2
	IDENTIFIER   = 3
	INT_CONST    = 4
	STRING_CONST = 5
)

// Constants for keywords
const (
	CLASS       = 10
	METHOD      = 11
	FUNCTION    = 12
	CONSTRUCTOR = 13
	INT         = 14
	BOOLEAN     = 15
	CHAR        = 16
	VOID        = 17
	VAR         = 18
	STATIC      = 19
	FIELD       = 20
	LET         = 21
	DO          = 22
	IF          = 23
	ELSE        = 24
	WHILE       = 25
	RETURN      = 26
	TRUE        = 27
	FALSE       = 28
	NULL        = 29
	THIS        = 30
)

type JackTokenizer struct {
	file         *os.File
	scanner      *bufio.Scanner
	currentToken string
}

func NewJackTokenizer(inputFile string) *JackTokenizer {
	// fill in here
}

func (jt *JackTokenizer) HasMoreTokens() bool {
	// fill in here
}

func (jt *JackTokenizer) Advance() {
	// fill in here
}

func (jt *JackTokenizer) TokenType() int {
	// fill in here
}

func (jt *JackTokenizer) KeyWord() int {
	// fill in heres
}

func (jt *JackTokenizer) Symbol() byte {
	// fill in here
}

func (jt *JackTokenizer) Identifier() string {
	// fill in here
}

func (jt *JackTokenizer) IntVal() int {
	// fill in here
}

func (jt *JackTokenizer) StringVal() string {
	// fill in here
}
