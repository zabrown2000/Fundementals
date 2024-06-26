package tokeniser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Constants for token types
const (
	KEYWORD      = 1
	SYMBOL       = 2
	IDENTIFIER   = 3
	INT_CONST    = 4
	STRING_CONST = 5
)

type Token struct {
	Token_type    int
	Token_content string
}

type Tokeniser struct {
	file         *os.File
	currentLine  string
	reader       *bufio.Reader
	Tokens       []Token
	LengthTokens int
	keywordsMap  map[string]bool
	symbolsMap   map[string]bool
}

// New creates a new tokeniser instance for a given file path
func New(path string) (*Tokeniser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	t := &Tokeniser{file: f, reader: reader}
	//t.Advance() // Advance to the first line
	t.init()
	return t, nil
}

// HasMoreLines returns true if there are more lines to be parsed
func (t *Tokeniser) HasMoreLines() bool {
	_, err := t.reader.Peek(1) // Peek to check for more lines without advancing
	return err == nil
}

// parseNextLine parses the next line of text, removing comments, whitespace, and empty lines
func (t *Tokeniser) parseNextLine() string {
	line, err := t.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if len(line) > 0 {
				// Handle the last line without a newline character
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "//") {
					line = ""
				}
				return line
			}
			return "" // Indicate no more lines

		}
		panic(fmt.Sprintf("err - couldn't get a line! %v", err))
	}
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "//") {
		if t.HasMoreLines() {
			line = t.parseNextLine()
		}
	}
	return line
}

func (t *Tokeniser) init() {
	keywords := []string{"class", "constructor", "method", "function", "field", "static", "var", "int", "char", "boolean", "void", "true", "false", "null",
		"this", "let", "do", "if", "else", "while", "return"}
	symbols := []string{"{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"}
	//Map for keywords
	t.keywordsMap = make(map[string]bool)
	for _, v := range keywords {
		t.keywordsMap[v] = true
	}

	t.symbolsMap = make(map[string]bool)
	for _, v := range symbols {
		t.symbolsMap[v] = true
	}
}

func (t *Tokeniser) Identifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z_]`)
	if !re.MatchString(s) {
		return false
	}
	for i, r := range s {
		if i == 0 {
			continue
		}

		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}

	}
	return true
}

func (t *Tokeniser) IntVal(s string) bool {
	num, err := strconv.Atoi(s)
	if err != nil {
		return false // Conversion failed, not a valid integer
	}

	// Check if the integer is within the specified range [0, 32767]
	return num >= 0 && num <= 32767
}

func (t *Tokeniser) StringVal(s string) bool {
	return len(s) > 1 && strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

func (t *Tokeniser) BlockComment(s string) bool {
	return len(s) > 3 && strings.HasPrefix(s, "/*") && strings.HasSuffix(s, "*/")
}

func (t *Tokeniser) LineComment(s string) bool {
	return len(s) > 2 && strings.HasPrefix(s, "//")
}

func (t *Tokeniser) TokeniseFile() {
	//fmt.Println("in tokeniseFile")
	inBlockComment := false

	for {

		line := t.parseNextLine() // parser returns line by line
		// Handle block comments
		if inBlockComment {
			if endIdx := strings.Index(line, "*/"); endIdx != -1 {
				line = line[endIdx+2:]
				inBlockComment = false
			} else {
				continue
			}
		}

		// Check for the start of a block comment
		if startIdx := strings.Index(line, "/*"); startIdx != -1 {
			endIdx := strings.Index(line[startIdx:], "*/")
			if endIdx != -1 {
				line = line[:startIdx] + line[startIdx+endIdx+2:]
			} else {
				line = line[:startIdx]
				inBlockComment = true
			}
		}
		chars := []rune(line)
		cur_char := ""
		token_candidate := ""

		for i := 0; i <= len(chars); i++ {
			if i < len(chars) {
				cur_char = string(chars[i])
			}

			// Check for line comment start and skip the rest of the line
			if i < len(chars)-1 && cur_char == "/" && string(chars[i+1]) == "/" {
				break
			}

			if t.keywordsMap[token_candidate] && !t.Identifier(cur_char) {
				t.Tokens = append(t.Tokens, Token{KEYWORD, token_candidate})
				token_candidate = ""
				t.LengthTokens++
			} else if t.symbolsMap[token_candidate] && (token_candidate != "/" || cur_char != "*") {
				t.Tokens = append(t.Tokens, Token{SYMBOL, token_candidate})
				token_candidate = ""
				t.LengthTokens++
			} else if t.StringVal(token_candidate) {
				t.Tokens = append(t.Tokens, Token{STRING_CONST, token_candidate[1 : len(token_candidate)-1]})
				token_candidate = ""
				t.LengthTokens++
			} else if t.Identifier(token_candidate) && !((cur_char >= "a" && cur_char <= "z") || (cur_char >= "A" && cur_char <= "Z") || (cur_char >= "0" && cur_char <= "9") || cur_char == "_") {
				t.Tokens = append(t.Tokens, Token{IDENTIFIER, token_candidate})
				token_candidate = ""
				t.LengthTokens++
			} else if t.IntVal(token_candidate) && !(cur_char >= "0" && cur_char <= "9") {
				t.Tokens = append(t.Tokens, Token{INT_CONST, token_candidate})
				token_candidate = ""
				t.LengthTokens++
			} else if t.BlockComment(token_candidate) {
				token_candidate = ""
				t.LengthTokens++

			}

			if cur_char == "\n" {
			} else if cur_char == " " && !strings.HasPrefix(token_candidate, "\"") {
			} else {
				token_candidate += cur_char
			}
		}

		if !t.HasMoreLines() {
			break
		}
	}
}
