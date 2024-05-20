package parser

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Tali: ctor, hasMoreCommands, advance, commandType, ar1, arg2
// Question - test files have lines that are started with '//' how to ignore them?

type CommandType int

//constructor

type Parser struct {
	file        *os.File
	scanner     *bufio.Scanner
	currentLine string
	hasMore     bool
	reader      *bufio.Reader
}

// New creates a new Parser instance for a given file path
func New(path string) (*Parser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	scanner := bufio.NewScanner(reader)
	p := &Parser{file: f, scanner: scanner}
	p.Advance() // Advance to the first line
	return p, nil
}

// advance moves to the next line and updates currentLine and hasMore
func (p *Parser) Advance() {
	p.parseNextLine()
}

// HasMoreLines returns true if there are more lines to be parsed
func (p *Parser) HasMoreLines() bool {
	return p.scanner.Scan()
}

// parseNextLine parses the next line of text, removing comments, whitespace, and empty lines
func (p *Parser) parseNextLine() {
	for {
		line, err := p.reader.ReadString('\n')
		if err != nil {
			break
		}
		if len(line) == 0 {
			if !p.HasMoreLines() {
				break // Exit the loop if there are no more lines
			}
			p.Advance() // Move to the next line
			continue    // Skip empty lines
		}
		if comment := strings.Index(line, "//"); comment > -1 {
			line = strings.TrimSpace(line[:comment])
		} else {
			line = strings.TrimSpace(line)
		}
		p.currentLine = line
		return // Exit the function after updating currentLine

	}
}

/*
// Parse returns the current line text content, excluding white spaces and comments
func (p *Parser) Parse() string {
	line := p.currentLine
	p.advance() // Move to the next line
	return line
}
*/

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
	C_IF                            // if
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
	case "if":
		return C_IF
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

// arg1 returns the first word of currentLine if it is of type arithmetic, otherwise it returns the second word.
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

// arg2 returns the third word of currentLine if it exists, otherwise it returns an empty string.
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
