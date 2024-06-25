package compilationEngine

import (
	"Fundementals/JackTranslator/tokeniser"
	"bufio"
	"os"
)

// TO DO: once all works, put checking type into its own function - int, bool, char (keywords), or identifier

type CompilationEngine struct {
	tokeniser         *tokeniser.Tokeniser
	plainWriter       *bufio.Writer
	hierarchWriter    *bufio.Writer
	currentToken      *tokeniser.Token
	currentTokenIndex int
}

func New(plainOutputFile string, hierarchOutputFile string, tokeniser *tokeniser.Tokeniser) *CompilationEngine {
	plainFile, err := os.OpenFile(plainOutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil
	}
	plainWriter := bufio.NewWriter(plainFile)

	hierarchFile, err := os.OpenFile(hierarchOutputFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return nil
	}
	hierarchWriter := bufio.NewWriter(hierarchFile)

	return &CompilationEngine{
		tokeniser:         tokeniser,
		plainWriter:       plainWriter,
		hierarchWriter:    hierarchWriter,
		currentTokenIndex: 0,
	}
}

// CompileClass compiles a complete class.
func (ce *CompilationEngine) CompileClass() {
	// 'class' identifier '{' classVarDec* subroutineDec* '}'

	// Write the opening tag <class>.
	ce.WriteOpenTag(ce.hierarchWriter, "class")
	ce.WriteOpenTag(ce.plainWriter, "tokens")
	// Advance the tokeniser to the next token and expect the keyword class.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "class") { // not a keyword and/or not 'class'
		panic("Unexpected token type! Expected keyword class")
	}
	// Write the class keyword.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier")
	}
	// Write the class name.
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Advance the tokeniser and expect the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") { // not a symbol and/or not '{'
		panic("Unexpected token! Expected {")
	}
	// Write the { symbol.
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Loop to handle class variable declarations (static or field) and subroutine declarations (constructor, function, or method):
	for { // TC before now, you had getToken before all comparisons of content - this time not - intentional? if not, outside or in the for loop?
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "static" || ce.currentToken.Token_content == "field") {
			// TC this is for inside CompileClassVarDec - but you didn't write it to static or field to file before advancing
			//so when you go to write you lost it - moving to inside CompileClassVarDec
			//ce.GetToken()
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
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </class>.
	ce.WriteCloseTag(ce.hierarchWriter, "class")
	ce.WriteCloseTag(ce.plainWriter, "tokens")
}

// CompileClassVarDec compiles a static declaration or a field declaration.
func (ce *CompilationEngine) CompileClassVarDec() {
	// ('static'|'field') type identifier (',' identifier)* ';'

	// Write the opening tag <classVarDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "classVarDec")
	// Write the current token (static or field).
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content) //no need to check if field or static, did that in compileClass
	// Advance the tokeniser and write the type (e.g., int, boolean, char, or a class name).
	ce.GetToken()
	// TC technically we should be calling CompileType and CompileClassName not just checking here
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for var type or identifier")
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	}
	// Advance the tokeniser and write the first variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	// Write the semicolon ; symbol.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token type! Expected symbol ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ";")
	ce.WriteXML(ce.plainWriter, "symbol", ";")
	//7. Write the closing tag </classVarDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "classVarDec")
}

// CompileSubroutine compiles a complete method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutine() {
	// ('constructor'|'function'|'method') ('void'|type) identifier '(' parameterList ')' subroutineBody

	// Write the opening tag <subroutineDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "subroutineDec")
	// Write the current token (constructor, function, or method).
	if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method") {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	// Advance the tokeniser and write the return type (void or a type).
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "void" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or a type
		panic("Unexpected token type! Expected keyword or identifier for subroutine return type")
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	}
	// Advance the tokeniser and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected identifier for subroutine name")
	}

	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token! Expected (")
	}
	// Advance the tokeniser and call compileParameterList.
	ce.GetToken()
	if ce.currentToken.Token_content != ")" {
		ce.CompileParameterList()
	}
	// Write the closing parenthesis ) - when no parameters token from getToken above, otherwise from getToken in broken loop in CompileParameterList
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected symbol )")
	}
	// Advance the tokeniser and call compileSubroutineBody.
	// TC for sake of uniformity - call getToken from inside CompileSubroutineBody? - already called in CompileParameterList -
	// TC this will advance twice before we check anything so removing here
	// ZB just like with paramList we do the intial getToken outside of the func and then more inside
	// ZB we need this before we call because the currentToken at this point will be ) and we need to move forward
	ce.GetToken()
	ce.CompileSubroutineBody()
	// Write the closing tag </subroutineDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineDec")
}

// CompileParameterList compiles a parameter list.
func (ce *CompilationEngine) CompileParameterList() {
	// ((type identifier) (',' type identifier)*)?

	// Write the opening tag <parameterList>
	ce.WriteOpenTag(ce.hierarchWriter, "parameterList")
	// Loop to handle parameters:
	for {
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
			ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
			panic("Unexpected token type! Expected keyword for var type or identifier")
		}
		if ce.currentToken.Token_type == tokeniser.KEYWORD {
			ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		}
		// Advance the tokeniser and write the variable name.
		ce.GetToken()
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name")
		}
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
		} else {
			break //no more params
		}
	}
	// Write the closing tag </parameterList>.
	ce.WriteCloseTag(ce.hierarchWriter, "parameterList")
}

// CompileSubroutineBody compiles the body of a method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutineBody() {
	// '{' varDec* statements '}'

	// Write the opening tag <subroutineBody>.
	ce.WriteOpenTag(ce.hierarchWriter, "subroutineBody")
	// Write the opening brace {.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Loop to handle variable declarations (var):
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "var" {
			ce.CompileVarDec()
		} else {
			break //no more vars
		}
	}
	// Call compileStatements to handle the statements within the subroutine body.
	ce.CompileStatements()
	// Write the closing brace }. - break after getting next token
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </subroutineBody>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineBody")
}

// CompileVarDec compiles a var declaration.
func (ce *CompilationEngine) CompileVarDec() {
	// 'var' type identifier (',' identifier)* ';'

	// Write the opening tag <varDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "varDec")
	// Write the current token var. - did check for it in caller function
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the type.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for var type or identifier")
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	}
	// Advance the tokeniser and write the first variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	// Write the semicolon ;.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected symbol ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </varDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "varDec")
}

// CompileStatements compiles a sequence of statements.
func (ce *CompilationEngine) CompileStatements() {
	// (letStatement|ifStatement|whileStatement|doStatement|returnStatement)*

	// Write the opening tag <statements>.
	ce.WriteOpenTag(ce.hierarchWriter, "statements")
	// Loop to handle each statement:
	for {
		ce.GetToken()
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
	}
	// Write the closing tag </statements>
	ce.WriteCloseTag(ce.hierarchWriter, "statements")
}

// CompileLet compiles a let statement.
func (ce *CompilationEngine) CompileLet() {
	// 'let' identifier ('[' expression ']')? '=' expression ';'

	// Write the opening tag <letStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "letStatement")
	// Write the current token let.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for var name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	// Advance the tokeniser to check for array indexing: If the current token is an opening bracket [,
	// write the bracket and call compileExpression. Write the closing bracket ] and advance the tokeniser.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.GetToken()
		ce.CompileExpression()
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]" {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
		}
	}
	// Write the equals sign =.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "=") {
		panic("Unexpected token! Expected =")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the semicolon ;.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </letStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "letStatement")
}

// CompileDo compiles a do statement.
func (ce *CompilationEngine) CompileDo() {
	// 'do' subroutineCall ';'

	// Write the opening tag <doStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "doStatement")
	// Write the current token do.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileSubroutineCall
	ce.GetToken()
	ce.CompileSubroutineCall()
	// Advance the tokeniser and write the semicolon ;.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </doStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "doStatement")
}

// CompileWhile compiles a while statement.
func (ce *CompilationEngine) CompileWhile() {
	// 'while' '('expression')' '{' statements '}'

	// Write the opening tag <whileStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "whileStatement")
	// Write the current token while.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the closing parenthesis ).
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Call compileStatements.
	ce.CompileStatements()
	// Write the closing brace }.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </whileStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "whileStatement")
}

// CompileReturn compiles a return statement.
func (ce *CompilationEngine) CompileReturn() {
	// 'return' expression? ';'

	// Write the opening tag <returnStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "returnStatement")
	// Write the current token return.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser to check for an expression:
	//     If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		ce.CompileExpression()
	}
	// Write the semicolon ;. - get token before leaving compileExpression
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Write the closing tag </returnStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "returnStatement")
}

// CompileIf compiles an if statement, possibly with a trailing else clause.
func (ce *CompilationEngine) CompileIf() {
	// 'if' '('expression')' '{'statements'}' ('else' '{'statements'}')?

	// Write the opening tag <ifStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "ifStatement")
	// Write the current token if.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// Write the closing parenthesis ). - get token before leave compileExpression
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Call compileStatements.
	ce.CompileStatements()
	//8. Write the closing brace }.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	// Advance the tokeniser to check for an else clause: If the current token is else, write the keyword else,
	// the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "else" {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
			panic("Unexpected token! Expected {")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.CompileStatements()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
			panic("Unexpected token! Expected }")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		ce.GoBackToken()
	}
	// Write the closing tag </ifStatement>
	ce.WriteCloseTag(ce.hierarchWriter, "ifStatement")
}

// CompileExpression compiles an expression.
func (ce *CompilationEngine) CompileExpression() {
	// term (op term)*

	// Write the opening tag <expression>.
	ce.WriteOpenTag(ce.hierarchWriter, "expression")
	// Call compileTerm. --got token before calling compileExpression
	ce.CompileTerm()
	//3. Loop to handle additional terms connected by operators:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && (ce.currentToken.Token_content == "+" || ce.currentToken.Token_content == "-" || ce.currentToken.Token_content == "*" || ce.currentToken.Token_content == "/" || ce.currentToken.Token_content == "&" || ce.currentToken.Token_content == "|" || ce.currentToken.Token_content == "<" || ce.currentToken.Token_content == ">" || ce.currentToken.Token_content == "=") {
			str := ce.currentToken.Token_content
			if str == "<" {
				str = "&lt;"
			} else if str == ">" {
				str = "&gt;"
			} else if str == "&" {
				str = "&amp;"
			} else if str == `"` {
				str = "&quote;"
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", str)
			ce.WriteXML(ce.plainWriter, "symbol", str)
			ce.GetToken()
			ce.CompileTerm()
		} else {
			break
		}
	}
	// Write the closing tag </expression>
	ce.WriteCloseTag(ce.hierarchWriter, "expression")
}

// CompileTerm compiles a term.
func (ce *CompilationEngine) CompileTerm() {
	// integerConstant|stringConstant|keywordConstant|identifier|identifier'['expression']'|subroutineCall|
	// '(' expression ')'|unaryOp term

	// Write the opening tag <term>.
	ce.WriteOpenTag(ce.hierarchWriter, "term")
	// Depending on the current token, handle different types of terms:
	if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "true" || ce.currentToken.Token_content == "false" || ce.currentToken.Token_content == "null" || ce.currentToken.Token_content == "this") { // got token before called compileTerm
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		// save current id, and then get next to see if symbol - look ahead
		identifier := ce.currentToken.Token_content
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
			// array
			ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
			ce.WriteXML(ce.plainWriter, "identifier", identifier) // arrayName
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // [
			ce.GetToken()
			ce.CompileExpression() // end with break after getToken
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
				panic("Unexpected token! Expected ]")
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // [
		} else if ce.currentToken.Token_type == tokeniser.SYMBOL && (ce.currentToken.Token_content == "(" || ce.currentToken.Token_content == ".") {
			// subroutine call
			// need to go back so function is at subroutine name
			ce.GoBackToken()
			ce.CompileSubroutineCall()
		} else {
			ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
			ce.WriteXML(ce.plainWriter, "identifier", identifier)
			// need to move token back one since not using current token here
			ce.GoBackToken()
		}
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL {
		symbol := ce.currentToken.Token_content
		if symbol == "(" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileExpression()
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
				panic("Unexpected token! Expected )")
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // )
		} else if symbol == "-" || symbol == "~" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileTerm()
		}
	} else if ce.currentToken.Token_type == tokeniser.INT_CONST {
		ce.WriteXML(ce.hierarchWriter, "int_const", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "int_const", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == tokeniser.STRING_CONST {
		ce.WriteXML(ce.hierarchWriter, "string_const", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "string_const", ce.currentToken.Token_content)
	}
	// Write the closing tag </term>.
	ce.WriteCloseTag(ce.hierarchWriter, "term")
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
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
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
