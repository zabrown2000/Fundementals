package compilationEngine

import (
	"Fundementals/JackTranslator/symbolTable"
	"Fundementals/JackTranslator/tokeniser"
	"Fundementals/JackTranslator/vmWriter"
	"bufio"
)

/*
This class does the compilation itself. It reads its input from a JackTokenizer and writes its output into a
VMWriter. It is organized as a series of compilexxx ( ) routines, where xxx is a syntactic element of the
Jack language. The contract between these routines is that each compilexxx ( ) routine should read the
syntactic construct xxx from the input, advance ( ) the tokenizer exactly beyond xxx, and emit to the output
VM code effecting the semantics of xxx. Thus compilexxx ( ) may only be called if indeed xxx is the next
syntactic element of the input. If xxx is a part of an expression and thus has a value, the emitted code
should compute this value and leave it at the top of the VM stack.
The API of this module is identical to that of the syntax analyzer’s compilation-Engine module fromchapter 10,
and thus we suggest gradually morphing the syntax analyzer into a full compiler. Section 11.5
provides step-by-step instructions and test programs for this construction.
*/

// TO DO: refactor checking symbols to be own func, like send symbol wanted and panic with currentToken and expected symbol

type CompilationEngine struct {
	tokeniser *tokeniser.Tokeniser
	//plainWriter       *bufio.Writer
	//hierarchWriter    *bufio.Writer
	vmWriter              *vmWriter.VMWriter
	symbolTable           *symbolTable.SymbolTable
	currentClassName      string
	currentSubroutineName string
	currentToken          *tokeniser.Token
	currentTokenIndex     int
	labelIndex            int
}

func New(outputFile string, tokeniser *tokeniser.Tokeniser) *CompilationEngine {
	//plainFile, err := os.OpenFile(plainOutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	//if err != nil {
	//	return nil
	//}
	//plainWriter := bufio.NewWriter(plainFile)
	//hierarchFile, err := os.OpenFile(hierarchOutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	//if err != nil {
	//	return nil
	//}
	//hierarchWriter := bufio.NewWriter(hierarchFile)

	return &CompilationEngine{
		tokeniser:         tokeniser,
		vmWriter:          vmWriter.New(outputFile),
		symbolTable:       symbolTable.New(),
		currentTokenIndex: 0,
		labelIndex:        0,
	}
}

// CompileClass compiles a complete class.
func (ce *CompilationEngine) CompileClass() {
	// 'class' identifier '{' classVarDec* subroutineDec* '}'

	// Advance the tokeniser to the next token and expect the keyword class.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "class") { // not a keyword and/or not 'class'
		panic("Unexpected token type! Expected keyword class")
	}
	// Write the class keyword.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for className")
	}
	ce.currentClassName = ce.currentToken.Token_content
	// Write the class name.
	//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Advance the tokeniser and expect the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") { // not a symbol and/or not '{'
		panic("Unexpected token! Expected {")
	}
	// Write the { symbol.
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Loop to handle class variable declarations (static or field) and subroutine declarations (constructor, function, or method):
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "static" || ce.currentToken.Token_content == "field") {
			ce.CompileClassVarDec()
		} else if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method") {
			ce.CompileSubroutine()
		} else {
			break // class is complete
		}
	}
	// Write the closing brace } symbol.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") { // not a symbol and/or not '}'
		panic("Unexpected token! Expected }")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileClassVarDec compiles a static declaration or a field declaration.
func (ce *CompilationEngine) CompileClassVarDec() {
	// ('static'|'field') type identifier (',' identifier)* ';'

	// Write the current token (static or field).
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content) //no need to check if field or static, did that in compileClass
	varKind := ce.currentToken.Token_content // getting kind for symbol table
	// Advance the tokeniser and write the type
	ce.GetToken()
	//compile the type
	varType := ce.currentToken.Token_content
	ce.CompileType()
	// Write the first variable name.
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	varName := ce.currentToken.Token_content // get first var name for symbol table
	//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	ce.symbolTable.Define(varName, varType, varKind)
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			varName = ce.currentToken.Token_content // get var name for symbol table
			ce.symbolTable.Define(varName, varType, varKind)
			//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	// Write the semicolon ; symbol.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token type! Expected symbol ;")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ";")
	//ce.WriteXML(ce.plainWriter, "symbol", ";")
}

// CompileSubroutine compiles a complete method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutine() {
	// ('constructor'|'function'|'method') ('void'|type) identifier '(' parameterList ')' subroutineBody

	// Write the current token (constructor, function, or method).
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method")) {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	if ce.currentToken.Token_content == "method" {
		ce.symbolTable.Define("this", ce.currentClassName, "argument")
	}

	// Advance the tokeniser and write the return type (void or a type).
	ce.GetToken()
	if !((ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "void" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "boolean")) || (ce.currentToken.Token_type == tokeniser.IDENTIFIER)) { // not a keyword or a type
		panic("Unexpected token type! Expected keyword or identifier for subroutine return type")
	}
	//if ce.currentToken.Token_type == tokeniser.KEYWORD {
	//	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
	//	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//}
	// Advance the tokeniser and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.currentSubroutineName = ce.currentToken.Token_content
		//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected identifier for subroutine name")
	}

	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	// Advance the tokeniser and call compileParameterList.
	ce.GetToken()
	ce.CompileParameterList()

	// Write the closing parenthesis ) - when no parameters token from getToken above, otherwise from getToken in broken loop in CompileParameterList
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token type! Expected symbol )")
	}
	// Advance the tokeniser and call compileSubroutineBody.
	ce.GetToken()
	ce.CompileSubroutineBody()
}

// CompileParameterList compiles a parameter list.
func (ce *CompilationEngine) CompileParameterList() {
	// ((type identifier) (',' type identifier)*)?

	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		return
	} else {
		ce.GoBackToken() //will get next token in loop
	}
	// Loop to handle parameters:
	for {
		ce.GetToken()
		//compile the type
		varType := ce.currentToken.Token_content
		ce.CompileType()
		// Write the variable name.
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name")
		}
		varName := ce.currentToken.Token_content
		ce.symbolTable.Define(varName, varType, "argument")
		//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ",") {
			break //no more params
		}
	}
}

// CompileSubroutineBody compiles the body of a method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutineBody() {
	// '{' varDec* statements '}'

	// Write the opening brace {.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Loop to handle variable declarations (var):
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "var" {
			ce.CompileVarDec()
		} else {
			break //no more vars
		}
	}
	// Handle bootstrap code for dif subroutine types
	ce.vmWriter.WriteFunction(ce.currentClassName+"."+ce.currentSubroutineName, ce.symbolTable.VarCount("local"))
	if ce.currentToken.Token_content == "constructor" {
		ce.vmWriter.WritePush("constant", ce.symbolTable.VarCount("field"))
		ce.vmWriter.WriteCall("Memory.alloc", 1)
		ce.vmWriter.WritePop("pointer", 0)
	} else if ce.currentToken.Token_content == "method" {
		ce.vmWriter.WritePush("argument", 0)
		ce.vmWriter.WritePop("pointer", 0)
	}
	// Call compileStatements to handle the statements within the subroutine body.
	ce.CompileStatements()
	// Write the closing brace }. - break after getting next token
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileVarDec compiles a var declaration.
func (ce *CompilationEngine) CompileVarDec() {
	// 'var' type identifier (',' identifier)* ';'

	// Write the current token var. - did check for it in caller function
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the type.
	ce.GetToken()
	//compile the type
	varType := ce.currentToken.Token_content
	ce.CompileType()
	// Write the first variable name.
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	varName := ce.currentToken.Token_content
	ce.symbolTable.Define(varName, varType, "local")
	//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			varName = ce.currentToken.Token_content
			ce.symbolTable.Define(varName, varType, "local")
			//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	// Write the semicolon ;.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected symbol ;")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileStatements compiles a sequence of statements.
func (ce *CompilationEngine) CompileStatements() {
	// (letStatement|ifStatement|whileStatement|doStatement|returnStatement)*

	// Loop to handle each statement:
	for {
		if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "let" {
			ce.CompileLet()
		} else if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "if" {
			ce.CompileIf()
		} else if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "while" {
			ce.CompileWhile()
		} else if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "do" {
			ce.CompileDo()
		} else if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "return" {
			ce.CompileReturn()
		} else {
			break
		}
		ce.GetToken()
	}
}

// CompileLet compiles a let statement.
func (ce *CompilationEngine) CompileLet() {
	// 'let' identifier ('[' expression ']')? '=' expression ';'

	// Write the current token let.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for var name")
	}
	varName := ce.currentToken.Token_content
	//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Advance the tokeniser to check for array indexing: If the current token is an opening bracket [,
	// write the bracket and call compileExpression. Write the closing bracket ] and advance the tokeniser.
	ce.GetToken()
	isArray := false
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
		//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		isArray = true
		ce.GetToken()
		ce.CompileExpression()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
			panic("Unexpected token! Expected ]")
		}
		ce.vmWriter.WritePush(ce.symbolTable.KindOf(varName), ce.symbolTable.IndexOf(varName))
		ce.vmWriter.WriteArithmetic("add")
		ce.GetToken()
	}
	// Write the equals sign =.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "=") {
		panic("Unexpected token! Expected =")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the semicolon ;.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	if isArray {
		ce.vmWriter.WritePop("temp", 0)
		ce.vmWriter.WritePop("pointer", 1)
		ce.vmWriter.WritePush("temp", 0)
		ce.vmWriter.WritePop("that", 0)
	} else {
		ce.vmWriter.WritePop(ce.symbolTable.KindOf(varName), ce.symbolTable.IndexOf(varName))
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileDo compiles a do statement.
func (ce *CompilationEngine) CompileDo() {
	// 'do' subroutineCall ';'

	// Write the current token do.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileSubroutineCall
	ce.GetToken()
	ce.CompileSubroutineCall()
	// Advance the tokeniser and write the semicolon ;.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	ce.vmWriter.WritePop("temp", 0)
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileWhile compiles a while statement.
func (ce *CompilationEngine) CompileWhile() {
	// 'while' '('expression')' '{' statements '}'

	// Write the current token while.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the closing parenthesis ).
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	// Write the closing brace }.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
}

// CompileReturn compiles a return statement.
func (ce *CompilationEngine) CompileReturn() {
	// 'return' expression? ';'

	// Write the current token return.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser to check for an expression: If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		ce.CompileExpression()
	} else {
		ce.vmWriter.WritePush("constant", 0)
	}
	// Write the semicolon ;. - get token before leaving compileExpression
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	ce.vmWriter.WriteReturn()
}

// CompileIf compiles an if statement, possibly with a trailing else clause.
func (ce *CompilationEngine) CompileIf() {
	// 'if' '('expression')' '{'statements'}' ('else' '{'statements'}')?

	// Write the current token if.
	//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the closing parenthesis ). - get token before leave compileExpression
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	// Write the closing brace }.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser to check for an else clause: If the current token is else, write the keyword else,
	// the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "else" {
		//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
			panic("Unexpected token! Expected {")
		}
		//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.GetToken()
		ce.CompileStatements()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
			panic("Unexpected token! Expected }")
		}
		//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		ce.GoBackToken()
	}
}

// CompileExpression compiles an expression.
func (ce *CompilationEngine) CompileExpression() {
	// term (op term)*

	// Call compileTerm. --got token before calling compileExpression
	ce.CompileTerm()
	// Loop to handle additional terms connected by operators:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && isOperator(ce.currentToken.Token_content) {
			op := ce.currentToken.Token_content
			//if str == "<" {
			//	str = "&lt;"
			//} else if str == ">" {
			//	str = "&gt;"
			//} else if str == "&" {
			//	str = "&amp;"
			//} else if str == `"` {
			//	str = "&quote;"
			//}
			//ce.WriteXML(ce.hierarchWriter, "symbol", str)
			//ce.WriteXML(ce.plainWriter, "symbol", str)
			ce.GetToken()
			ce.CompileTerm()
			ce.WriteArithmeticCommand(op) //wait until term is placed in stack before putting in op for postfix notation
		} else {
			break
		}
	}
}

// CompileTerm compiles a term.
func (ce *CompilationEngine) CompileTerm() {
	// integerConstant|stringConstant|keywordConstant|identifier|identifier'['expression']'|subroutineCall|
	// '(' expression ')'|unaryOp term

	// Depending on the current token, handle different types of terms:
	if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "true" || ce.currentToken.Token_content == "false" || ce.currentToken.Token_content == "null" || ce.currentToken.Token_content == "this") { // got token before called compileTerm
		switch ce.currentToken.Token_content {
		case "true":
			ce.vmWriter.WritePush("constant", 0)
			ce.vmWriter.WriteArithmetic("not")
		case "false", "null":
			ce.vmWriter.WritePush("constant", 0)
		case "this":
			ce.vmWriter.WritePush("pointer", 0)
		}
		//ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		// save current identifier, and then get next token to see if symbol - look ahead
		identifier := ce.currentToken.Token_content
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
			// array
			//ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
			//ce.WriteXML(ce.plainWriter, "identifier", identifier) // arrayName
			//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // [
			ce.GetToken()
			ce.CompileExpression() // end with break after getToken
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
				panic("Unexpected token! Expected ]")
			}
			ce.vmWriter.WritePush(ce.symbolTable.KindOf(identifier), ce.symbolTable.IndexOf(identifier))
			ce.vmWriter.WriteArithmetic("add")
			ce.vmWriter.WritePop("pointer", 1)
			ce.vmWriter.WritePush("that", 0)
			//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // [
		} else if ce.currentToken.Token_type == tokeniser.SYMBOL && (ce.currentToken.Token_content == "(" || ce.currentToken.Token_content == ".") {
			// subroutine call
			// need to go back so function is at subroutine name
			ce.GoBackToken()
			ce.CompileSubroutineCall()
		} else {
			//ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
			//ce.WriteXML(ce.plainWriter, "identifier", identifier)
			ce.vmWriter.WritePush(ce.symbolTable.KindOf(identifier), ce.symbolTable.IndexOf(identifier))
			// need to move token back one since not using current token here
			ce.GoBackToken()
		}
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL {
		symbol := ce.currentToken.Token_content
		if symbol == "(" {
			//ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			//ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileExpression()
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
				panic("Unexpected token! Expected )")
			}
			//ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			//ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // )
		} else if symbol == "-" || symbol == "~" {
			//ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			//ce.WriteXML(ce.plainWriter, "symbol", symbol)
			unaryOp := ce.currentToken.Token_content
			ce.GetToken()
			ce.CompileTerm()
			if unaryOp == "-" {
				ce.vmWriter.WriteArithmetic("neg")
			} else if unaryOp == "~" {
				ce.vmWriter.WriteArithmetic("not")
			}
		}
	} else if ce.currentToken.Token_type == tokeniser.INT_CONST {
		ce.vmWriter.WritePush("constant", ce.currentToken.Token_content)
		//ce.WriteXML(ce.hierarchWriter, "integerConstant", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "integerConstant", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.STRING_CONST {
		stringVal := ce.currentToken.Token_content
		ce.vmWriter.WritePush("constant", len(stringVal))
		ce.vmWriter.WriteCall("String.new", 1) //1 param, length of string
		for _, char := range stringVal {       //NOTE: needs to be rune here?
			ce.vmWriter.WritePush("constant", int(char))
			ce.vmWriter.WriteCall("String.appendChar", 2) //above push and result of new strign will be in stack
		}
		//ce.WriteXML(ce.hierarchWriter, "stringConstant", ce.currentToken.Token_content)
		//ce.WriteXML(ce.plainWriter, "stringConstant", ce.currentToken.Token_content)
	}
}

// CompileExpressionList is responsible for compiling a (possibly empty) comma-separated list of expressions. This list is typically found within the argument list of a subroutine call.
func (ce *CompilationEngine) CompileExpressionList() {
	// (expression (',' expression)*)?

	// Write the opening tag <expressionList>.
	ce.WriteOpenTag(ce.hierarchWriter, "expressionList")
	// Check if the current token indicates the start of an expression. This can be identified by looking for tokens that can
	// start an expression such as integer constants, string constants, keyword constants, variable names, subroutine calls,
	// expressions enclosed in parentheses, and unary operators.
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		ce.WriteCloseTag(ce.hierarchWriter, "expressionList")
		return // empty list
	}
	// If there is at least one expression: Call compileExpression to compile the first expression.
	// Loop to handle additional expressions separated by commas
	ce.CompileExpression() // do getToken before break

	// while have comma, another expression
	for ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // ','
		ce.GetToken()
		ce.CompileExpression() // did a getToken before breaking
	}
	// Write the closing tag </expressionList>.
	ce.WriteCloseTag(ce.hierarchWriter, "expressionList")
}

// CompileSubroutineCall compiles a subroutine call
func (ce *CompilationEngine) CompileSubroutineCall() {
	// identifier'('expressionList')' | identifier'.'identifier'('expressionList')'

	// Write the subroutine name
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	//ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	//ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		// Advance the tokeniser and call compileExpressionList.
		ce.GetToken()
		ce.CompileExpressionList()
		// Write the closing parenthesis ). - when leave expressionList did a get token
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
			panic("Unexpected token! Expected )")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "." {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.GetToken()
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for class or var name")
		}
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(" {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			// Advance the tokeniser and call compileExpressionList.
			ce.GetToken()
			ce.CompileExpressionList()
			// Write the closing parenthesis ). - when leave expressionList did a get token
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
				panic("Unexpected token! Expected )")
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		}
	} else {
		panic("Unexpected token! Expected ( or .")
	}
}

func (ce *CompilationEngine) CompileType() {
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) && !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for type or identifier")
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		ce.GetToken()
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		ce.GetToken()
	}
}

// helper functions

func (ce *CompilationEngine) WriteXML(writer *bufio.Writer, tag string, content string) {
	toWrite := "<" + tag + ">" + content + "</" + tag + ">\n"
	_, err := writer.Write([]byte(toWrite))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func (ce *CompilationEngine) WriteOpenTag(writer *bufio.Writer, tag string) {
	toWrite := "<" + tag + ">\n"
	_, err := writer.Write([]byte(toWrite))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func (ce *CompilationEngine) WriteCloseTag(writer *bufio.Writer, tag string) {
	toWrite := "</" + tag + ">\n"
	_, err := writer.Write([]byte(toWrite))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func (ce *CompilationEngine) GetToken() {
	ce.currentToken = &ce.tokeniser.Tokens[ce.currentTokenIndex]
	ce.currentTokenIndex = ce.currentTokenIndex + 1
}

func (ce *CompilationEngine) GoBackToken() {
	ce.currentTokenIndex = ce.currentTokenIndex - 1
	ce.currentToken = &ce.tokeniser.Tokens[ce.currentTokenIndex-1]
}

func (ce *CompilationEngine) currentFunction() string {
	if len(ce.currentClassName) != 0 && len(ce.currentSubroutineName) != 0 {
		return ce.currentClassName + "." + ce.currentSubroutineName
	}
	return ""
}

func isOperator(symbol string) bool {
	return symbol == "+" || symbol == "-" || symbol == "*" || symbol == "/" ||
		symbol == "&" || symbol == "|" || symbol == "<" || symbol == ">" || symbol == "="
}

func (ce *CompilationEngine) WriteArithmeticCommand(command string) {
	switch command {
	case "+":
		ce.vmWriter.WriteArithmetic("add")
	case "-":
		ce.vmWriter.WriteArithmetic("sub")
	case "*":
		ce.vmWriter.WriteCall("Math.multiply", 2)
	case "/":
		ce.vmWriter.WriteCall("Math.divide", 2)
	case "&":
		ce.vmWriter.WriteArithmetic("and")
	case "|":
		ce.vmWriter.WriteArithmetic("or")
	case "<":
		ce.vmWriter.WriteArithmetic("lt")
	case ">":
		ce.vmWriter.WriteArithmetic("gt")
	case "=":
		ce.vmWriter.WriteArithmetic("eq")
	}
}
