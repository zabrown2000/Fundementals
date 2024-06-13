package JackTranslator

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var keywords []string
var symbols []string

// throws if illegal token
// use lexical elements section from slides

// Constants for token types
const (
	KEYWORD      = 1
	SYMBOL       = 2
	IDENTIFIER   = 3
	INT_CONST    = 4
	STRING_CONST = 5
)

type Token struct {
	token_type int
	token      string
}

type Tokeniser struct {
	file        *os.File
	currentLine string
	reader      *bufio.Reader
	tokens      *[]Token
}

// New creates a new tokeniser instance for a given file path
func New(path string) (*Tokeniser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	t := &Tokeniser{file: f, reader: reader}
	t.Advance() // Advance to the first line
	t.init()
	return t, nil
}

// Advance moves to the next line and updates currentLine and hasMore
func (t *Tokeniser) Advance() {
	t.parseNextLine()
}

// HasMoreLines returns true if there are more lines to be parsed
func (t *Tokeniser) HasMoreLines() bool {
	_, err := t.reader.Peek(1) // Peek to check for more lines without advancing
	return err == nil
}

// parseNextLine parses the next line of text, removing comments, whitespace, and empty lines
func (t *Tokeniser) parseNextLine() {
	line, err := t.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if len(line) > 0 {
				// Handle the last line without a newline character
				t.currentLine = strings.TrimSpace(line)
				return
			}
			t.currentLine = "" // Indicate no more lines
			return

		}
		panic(fmt.Sprintf("err - couldn't get a line! %v", err))
	}
	t.currentLine = line
	return // Exit the function after updating currentLine
}
func (t *Tokeniser) init() {
	keywords = []string{"class", "constructor", "method", "function", "field", "static", "var", "int", "char", "boolean", "void", "true", "false", "null",
		"this", "let", "do", "if", "else", "while", "return"}
	symbols = []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
	//Map for keywords
	keywordsMap := make(map[string]bool)
	for _, v := range keywords {
		keywordsMap[v] = true
	}
	symbolsMap := make(map[string]bool)
	for _, v := range symbols {
		symbolsMap[v] = true
	}
}

func (t *Tokeniser) HasMoreTokens() bool {
	return false // fill in here
}

func (t *Tokeniser) AdvanceToken() {
	// fill in here
}

func (t *Tokeniser) TokenType() int {
	return 0 // fill in here
}

func (t *Tokeniser) KeyWord() int {
	return 0 // fill in here
}

func (t *Tokeniser) Symbol() byte {
	return 0 // fill in here
}

func (t *Tokeniser) Identifier() string {
	return "" // fill in here
}

func (t *Tokeniser) IntVal() int {
	return 0 // fill in here
}

func (t *Tokeniser) StringVal() string {
	return "" // fill in here
}
