package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type CommandType int

type Parser struct {
	file        *os.File
	currentLine string
	reader      *bufio.Reader
}

// New creates a new Parser instance for a given file path
func New(path string) (*Parser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	p := &Parser{file: f, reader: reader}
	p.Advance() // Advance to the first line
	return p, nil
}

// Advance moves to the next line and updates currentLine and hasMore
func (p *Parser) Advance() {
	p.parseNextLine()
}

// HasMoreLines returns true if there are more lines to be parsed
func (p *Parser) HasMoreLines() bool {
	_, err := p.reader.Peek(1) // Peek to check for more lines without advancing
	return err == nil
}

// parseNextLine parses the next line of text, removing comments, whitespace, and empty lines
func (p *Parser) parseNextLine() {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			p.currentLine = "" // Indicate no more lines
			return
		}
		panic(fmt.Sprintf("err - couldn't get a line! %v", err))
	}
	p.currentLine = line
	return // Exit the function after updating currentLine
}

// Close closes the file being parsed
func (p *Parser) Close() error {
	return p.file.Close()
}

// Command types
const (
	C_ARITHMETIC CommandType = iota // arithmetic (add,sub,neg,eq,lt,gt,and,or,not)
	C_PUSH                          // push
	C_POP                           // pop
	C_LABEL                         // label
	C_GOTO                          // goto
	C_IFGOTO                        // if-goto
	C_FUNCTION                      // function
	C_RETURN                        // return
	C_CALL                          // call
)

func (p *Parser) CommandType() CommandType {
	line := p.currentLine
	// Split the line into tokens based on whitespace characters
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return -1 // Invalid command type
	}
	// Determine the command type based on the first token
	switch tokens[0] {
	case "add", "sub", "neg", "eq", "lt", "gt", "and", "or", "not":
		return C_ARITHMETIC
	case "push":
		return C_PUSH
	case "pop":
		return C_POP
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IFGOTO
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	case "call":
		return C_CALL
	default:
		return -1 // Invalid command type
	}
}

// Arg1 returns the first word of currentLine if it is of type arithmetic, otherwise it returns the second word.
func (p *Parser) Arg1() string {
	line := p.currentLine
	// Split the line into tokens based on whitespace characters
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return "" // Empty string if there are no tokens
	}
	if p.CommandType() == C_ARITHMETIC {
		return tokens[0] // First word for arithmetic commands
	} else if len(tokens) >= 2 {
		return tokens[1] // Second word for other commands
	}
	return "" // Empty string if there is no second word
}

// Arg2 returns the third word of currentLine if it exists, otherwise it returns an empty string.
func (p *Parser) Arg2() int {
	line := p.currentLine
	// Split the line into tokens based on whitespace characters
	tokens := strings.Fields(line)
	if len(tokens) < 3 {
		return -1 // Return -1 if there is no third word
	}
	// Parse the third word as an integer
	arg2, err := strconv.Atoi(tokens[2])
	if err != nil {
		return -1 // Return -1 if the third word cannot be parsed as an integer
	}
	return arg2
}
