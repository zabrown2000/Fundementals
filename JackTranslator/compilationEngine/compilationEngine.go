package compilationEngine

import (
	"Fundementals/JackTranslator/symbolTable"
	"Fundementals/JackTranslator/tokeniser"
	"Fundementals/JackTranslator/vmWriter"
	"strconv"
)

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
	// Advance the tokeniser and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for className")
	}
	ce.currentClassName = ce.currentToken.Token_content
	// Advance the tokeniser and expect the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") { // not a symbol and/or not '{'
		panic("Unexpected token! Expected {")
	}
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
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") { // not a symbol and/or not '}'
		panic("Unexpected token! Expected }")
	}
}

// CompileClassVarDec compiles a static declaration or a field declaration.
func (ce *CompilationEngine) CompileClassVarDec() {
	// ('static'|'field') type identifier (',' identifier)* ';'

	// getting kind for symbol table
	varKind := ce.currentToken.Token_content // will be static or field since class
	// Advance the tokeniser and write the type
	ce.GetToken()
	//compile the type
	varType := ce.currentToken.Token_content
	ce.CompileType()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	// getting type for symbol table
	varName := ce.currentToken.Token_content // get first var name for symbol table
	ce.symbolTable.Define(varName, varType, varKind)
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			varName = ce.currentToken.Token_content // get var name for symbol table
			ce.symbolTable.Define(varName, varType, varKind)
		} else {
			break // no more variables
		}
	}
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token type! Expected symbol ;")
	}
}

// CompileSubroutine compiles a complete method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutine() {
	// ('constructor'|'function'|'method') ('void'|type) identifier '(' parameterList ')' subroutineBody
	//we checked this before calling CompileSubroutine
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method")) {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	//if ce.currentToken.Token_content == "method" {
	//	ce.symbolTable.Define("this", ce.currentClassName, "argument")
	//}
	var isMethod bool
	if ce.currentToken.Token_content == "method" {
		isMethod = true
	} else {
		isMethod = false
	}
	ce.symbolTable.StartSubroutine(isMethod, ce.currentClassName)

	// Advance the tokeniser and write the return type (void or a type).
	ce.GetToken()
	if !((ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "void" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "boolean")) || (ce.currentToken.Token_type == tokeniser.IDENTIFIER)) { // not a keyword or a type
		panic("Unexpected token type! Expected keyword or identifier for subroutine return type")
	}
	// Advance the tokeniser and set the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.currentSubroutineName = ce.currentToken.Token_content
	} else {
		panic("Unexpected token type! Expected identifier for subroutine name")
	}

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
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ",") {
			break //no more params
		}
	}
}

// CompileSubroutineBody compiles the body of a method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutineBody() {
	// '{' varDec* statements '}'

	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
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
	// break after getting next token
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
}

// CompileVarDec compiles a var declaration.
func (ce *CompilationEngine) CompileVarDec() {
	// 'var' type identifier (',' identifier)* ';'

	// Advance the tokeniser
	ce.GetToken()
	//compile the type
	varType := ce.currentToken.Token_content //also assigning the type before we've verified if it's legal?
	ce.CompileType()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	varName := ce.currentToken.Token_content
	ce.symbolTable.Define(varName, varType, "local")
	// Loop to handle additional variables:
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.GetToken()
			if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			varName = ce.currentToken.Token_content
			ce.symbolTable.Define(varName, varType, "local")
		} else {
			break // no more variables
		}
	}
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected symbol ;")
	}
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

	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for var name")
	}
	varName := ce.currentToken.Token_content

	// Advance the tokeniser to check for array indexing: If the current token is an opening bracket [,
	// write the bracket and call compileExpression. Write the closing bracket ] and advance the tokeniser.
	ce.GetToken()
	isArray := false
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
		isArray = true
		ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(varName)), ce.symbolTable.IndexOf(varName))
		ce.GetToken()
		ce.CompileExpression()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
			panic("Unexpected token! Expected ]")
		}
		ce.vmWriter.WriteArithmetic("add")
		ce.GetToken()
	}
	// Write the equals sign =.
	// NOTE: might not need =?
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "=") {
		panic("Unexpected token! Expected =")
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//experimental addition of get token here
	if ce.currentToken.Token_content != ";" {
		ce.GetToken()
	}
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ; Let " + ce.currentToken.Token_content)
	}
	// add necessary commands if array
	if isArray {
		ce.vmWriter.WritePop("temp", 0)
		ce.vmWriter.WritePop("pointer", 1)
		ce.vmWriter.WritePush("temp", 0)
		ce.vmWriter.WritePop("that", 0)
	} else {
		ce.vmWriter.WritePop(ce.GetSeg(ce.symbolTable.KindOf(varName)), ce.symbolTable.IndexOf(varName))
	}
}

// CompileDo compiles a do statement.
func (ce *CompilationEngine) CompileDo() {
	// 'do' subroutineCall ';'

	// Advance the tokeniser and call compileSubroutineCall
	ce.GetToken()
	ce.CompileSubroutineCall()
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ; Do " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WritePop("temp", 0)
}

// CompileWhile compiles a while statement.
func (ce *CompilationEngine) CompileWhile() {
	// 'while' '('expression')' '{' statements '}'

	labelContinue := ce.NewLabel()
	labelTop := ce.NewLabel()
	ce.vmWriter.WriteLabel(labelTop)
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	ce.vmWriter.WriteArithmetic("not")
	ce.vmWriter.WriteIfGoto(labelContinue)
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	// Call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
}

// CompileReturn compiles a return statement.
func (ce *CompilationEngine) CompileReturn() {
	// 'return' expression? ';'

	// Advance the tokeniser to check for an expression: If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";" {
		ce.vmWriter.WritePush("constant", 0)
	} else {
		ce.CompileExpression()
		// Write the semicolon ;. - get token before leaving compileExpression
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
			panic("Unexpected token! Expected ; Return" + ce.currentToken.Token_content)
		}
	}
	ce.vmWriter.WriteReturn()
}

// CompileIf compiles an if statement, possibly with a trailing else clause.
func (ce *CompilationEngine) CompileIf() {
	// 'if' '('expression')' '{'statements'}' ('else' '{'statements'}')?

	labelElse := ce.NewLabel()
	labelEnd := ce.NewLabel()
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected (")
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	// get token before leave compileExpression
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected )")
	}
	ce.vmWriter.WriteArithmetic("not")
	ce.vmWriter.WriteIfGoto(labelElse)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected {")
	}
	// Call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	// Write the closing brace }.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected }")
	}
	ce.vmWriter.WriteGoTo(labelEnd)
	ce.vmWriter.WriteLabel(labelElse)
	// Advance the tokeniser to check for an else clause: If the current token is else, write the keyword else,
	// the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "else" {
		ce.GetToken()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
			panic("Unexpected token! Expected {")
		}
		ce.GetToken()
		ce.CompileStatements()
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
			panic("Unexpected token! Expected }")
		}
	} else {
		ce.GoBackToken()
	}
	ce.vmWriter.WriteLabel(labelEnd)
}

// CompileExpression compiles an expression.
func (ce *CompilationEngine) CompileExpression() {
	// term (op term)*

	// Call compileTerm. --got token before calling compileExpression
	ce.CompileTerm()
	//maybe check here that )? or add if not ) to for?
	//if ce.currentToken.Token_content == ")" {
	//	return
	//}
	// Loop to handle additional terms connected by operators:
	for {
		if ce.currentToken.Token_content != ")" {
			ce.GetToken()
		}
		if ce.currentToken.Token_type == tokeniser.SYMBOL && isOperator(ce.currentToken.Token_content) {
			op := ce.currentToken.Token_content
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
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		// save current identifier, and then get next token to see if symbol - look ahead
		identifier := ce.currentToken.Token_content
		ce.GetToken()
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
			// array
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(identifier)), ce.symbolTable.IndexOf(identifier))
			ce.GetToken()
			ce.CompileExpression() // end with break after getToken
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
				panic("Unexpected token! Expected ]")
			}
			ce.vmWriter.WriteArithmetic("add")
			ce.vmWriter.WritePop("pointer", 1)
			ce.vmWriter.WritePush("that", 0)
		} else if ce.currentToken.Token_type == tokeniser.SYMBOL && (ce.currentToken.Token_content == "(" || ce.currentToken.Token_content == ".") {
			// subroutine call
			// need to go back so function is at subroutine name
			ce.GoBackToken()
			ce.CompileSubroutineCall()
		} else {
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(identifier)), ce.symbolTable.IndexOf(identifier))
			// need to move token back one since not using current token here
			ce.GoBackToken()
		}
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL {
		symbol := ce.currentToken.Token_content
		if symbol == "(" {
			ce.GetToken()
			ce.CompileExpression()
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
				panic("Unexpected token! Expected )")
			}
		} else if symbol == "-" || symbol == "~" {
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
		int_token, err := strconv.Atoi(ce.currentToken.Token_content)
		if err != nil {
			return
		}
		ce.vmWriter.WritePush("constant", int_token)
	} else if ce.currentToken.Token_type == tokeniser.STRING_CONST {
		stringVal := ce.currentToken.Token_content
		ce.vmWriter.WritePush("constant", len(stringVal))
		ce.vmWriter.WriteCall("String.new", 1) //1 param, length of string
		for _, char := range stringVal {       //NOTE: needs to be rune here?
			ce.vmWriter.WritePush("constant", int(char))
			ce.vmWriter.WriteCall("String.appendChar", 2) //above push and result of new string will be in stack
		}
	}
}

// CompileExpressionList is responsible for compiling a (possibly empty) comma-separated list of expressions. This list is typically found within the argument list of a subroutine call.
func (ce *CompilationEngine) CompileExpressionList() int {
	// (expression (',' expression)*)?

	nArgs := 0
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		return nArgs // no args
	}
	// If there is at least one expression: Call compileExpression to compile the first expression.
	// Loop to handle additional expressions separated by commas
	ce.CompileExpression() // do getToken before break
	nArgs++

	// while have comma, another expression
	for ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
		ce.GetToken()
		ce.CompileExpression() // did a getToken before breaking
		nArgs++
	}
	return nArgs
}

// CompileSubroutineCall compiles a subroutine call
func (ce *CompilationEngine) CompileSubroutineCall() {
	// identifier'('expressionList')' | identifier'.'identifier'('expressionList')'

	//I don't think we need to do this - we've already checked that it's an identifier, and all we do is go
	//forward again with the token and again get either the '.' or '('

	// Write the subroutine name
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	name := ce.currentToken.Token_content
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	nArgs := 0
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(" {
		ce.vmWriter.WritePush("pointer", 0)
		name = ce.currentClassName + "." + name
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "." {
		ce.GetToken()
		objName := ce.currentToken.Token_content
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for class or var name")
		}
		//subroutineName := ce.currentToken.Token_content
		typ := ce.symbolTable.TypeOf(objName)
		if typ == "int" || typ == "boolean" || typ == "char" || typ == "void" {
			panic("Unexpected token type! Expected non-built in type")
		} else if typ == "" { //This or the symbol table type of need to be matched
			name += "." + objName
		} else {
			nArgs = 1
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(objName)), ce.symbolTable.IndexOf(objName))
			name = ce.symbolTable.TypeOf(objName) + "." + name
		}
		//ce.GetToken()
		//kind := ce.symbolTable.KindOf(name)
		//if kind == "NONE" {
		//	ce.vmWriter.WritePush(kind, ce.symbolTable.IndexOf(name))
		//	name = ce.symbolTable.TypeOf(name) + "." + subroutineName
		//} else {
		//	name = name + "." + subroutineName
		//}
	} else {
		panic("Unexpected token! Expected ( or .")
	}
	// Advance the tokeniser and call compileExpressionList.
	ce.GetToken()
	nArgs = ce.CompileExpressionList() //+ 1
	// Write the closing parenthesis ). - when leave expressionList did a get token
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected ) " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteCall(name, nArgs)
}

func (ce *CompilationEngine) CompileType() {
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) && !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for type or identifier")
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD {
		ce.GetToken()
	} else if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.GetToken()
	}
}

// helper functions

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

func (ce *CompilationEngine) NewLabel() string {
	var l = ce.labelIndex
	ce.labelIndex = l + 1
	return "LABEL_" + strconv.Itoa(l)
}

func (ce *CompilationEngine) GetSeg(kind string) string {
	if kind == "field" {
		return "this"
	} else if kind == "static" {
		return kind
	} else if kind == "local" {
		return kind
	} else if kind == "argument" {
		return kind
	} else {
		return ""
	}
}
