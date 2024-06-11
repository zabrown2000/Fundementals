package JackTranslator

import (
	"bufio"
)

type CompilationEngine struct {
	tokenizer *JackTokenizer
	writer    *bufio.Writer
}

func NewCompilationEngine(inputFile string, outputFile string) *CompilationEngine {
	// fill in here
}

func (ce *CompilationEngine) compileClass() {
	// fill in here
}

func (ce *CompilationEngine) compileClassVarDec() {
	// fill in here
}

func (ce *CompilationEngine) compileSubroutine() {
	// fill in here
}

func (ce *CompilationEngine) compileParameterList() {
	// fill in here
}

func (ce *CompilationEngine) compileVarDec() {
	// fill in here
}

func (ce *CompilationEngine) compileStatements() {
	// fill in here
}

func (ce *CompilationEngine) compileDo() {
	// fill in here
}

func (ce *CompilationEngine) compileLet() {
	// fill in here
}

func (ce *CompilationEngine) compileWhile() {
	// fill in here
}

func (ce *CompilationEngine) compileReturn() {
	// fill in here
}

func (ce *CompilationEngine) compileIf() {
	// fill in here
}

func (ce *CompilationEngine) compileExpression() {
	// fill in here
}

func (ce *CompilationEngine) compileTerm() {
	// fill in here
}

func (ce *CompilationEngine) compileExpressionList() {
	// fill in here
}
