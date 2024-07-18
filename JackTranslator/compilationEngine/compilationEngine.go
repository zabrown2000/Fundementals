package compilationEngine

import (
	"Fundementals/JackTranslator/symbolTable"
	"Fundementals/JackTranslator/tokeniser"
	"Fundementals/JackTranslator/vmWriter"
	"fmt"
	"runtime"
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
		panic("Unexpected token type! Expected keyword class " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for className " + ce.currentToken.Token_content)
	}
	ce.currentClassName = ce.currentToken.Token_content
	// Advance the tokeniser and expect the opening brace {.
	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") { // not a symbol and/or not '{'
		panic("Unexpected token! Expected { start of class " + ce.currentToken.Token_content)
	}
	ce.GetToken()
	if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "static" || ce.currentToken.Token_content == "field") {
		ce.CompileClassVarDec()
	}
	if ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method") {
		ce.CompileSubroutineDec()
	}
	//ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") { // not a symbol and/or not '}'
		panic("Unexpected token! Expected } end of class " + ce.currentToken.Token_content)
	}
}

// CompileClassVarDec compiles a static declaration or a field declaration.
func (ce *CompilationEngine) CompileClassVarDec() {
	// ('static'|'field') type identifier (',' identifier)* ';'
	if !(ce.currentToken.Token_content == "static" || ce.currentToken.Token_content == "field") {
		return //not a class var dec
	}
	// getting kind for symbol table
	varKind := ce.currentToken.Token_content // will be static or field since class
	// Advance the tokeniser and write the type
	ce.GetToken()
	ce.CompileType()
	//compile the type
	varType := ce.currentToken.Token_content
	// getting type for symbol table
	ce.GetToken()                                           //varName
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name class var dec " + ce.currentToken.Token_content)
	}
	// Loop to handle additional variables:
	for {
		varName := ce.currentToken.Token_content // get first var name for symbol table
		ce.symbolTable.Define(varName, varType, varKind)
		ce.GetToken() //comma or ;
		if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content != "," {
			break //no more varnames and cur token ; out of for
		}
		ce.GetToken()                                           //was a comma, now next varname
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name class var dec loop " + ce.currentToken.Token_content)
		}
	}
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token type! Expected symbol ; end class var dec " + ce.currentToken.Token_content)
	}
	ce.GetToken() //if next token field or static rinse and repeat, if not we leave the if in compile class and proceed to next if
	ce.CompileClassVarDec()

}

// CompileSubroutineDec compiles a complete method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutineDec() {
	// ('constructor'|'function'|'method') ('void'|type) identifier '(' parameterList ')' subroutineBody
	//we checked this before calling CompileSubroutine
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method")) {
		return
	}
	var isMethod bool
	if ce.currentToken.Token_content == "method" {
		isMethod = true
	} else {
		isMethod = false
	}
	ce.symbolTable.StartSubroutine(isMethod, ce.currentClassName, ce.currentToken.Token_content)
	// Advance the tokeniser and write the return type (void or a type).
	ce.GetToken() //int/void/cahr/boolean/identifier
	if !((ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "void" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "boolean")) || (ce.currentToken.Token_type == tokeniser.IDENTIFIER)) { // not a keyword or a type
		panic("Unexpected token type! Expected keyword or identifier for subroutine return type subroutdec " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and set the subroutine name.
	ce.GetToken() //subroutine name
	if ce.currentToken.Token_type == tokeniser.IDENTIFIER {
		ce.currentSubroutineName = ce.currentToken.Token_content
	} else {
		panic("Unexpected token type! Expected identifier for subroutine name subroutdec " + ce.currentToken.Token_content)
	}

	ce.GetToken()
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected ( start param list " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and call compileParameterList.
	ce.CompileParameterList()

	// Write the closing parenthesis ) - when no parameters token from getToken above, otherwise from getToken in broken loop in CompileParameterList
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token type! Expected symbol ) end param list " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and call compileSubroutineBody.
	ce.CompileSubroutineBody()
	ce.GetToken() // finished subroutine body with }, next token to see if more functions
	ce.CompileSubroutineDec()
}

// CompileParameterList compiles a parameter list.
func (ce *CompilationEngine) CompileParameterList() {
	// ((type identifier) (',' type identifier)*)?
	//enter with ( now so advance and check whether ) or not
	ce.GetToken() // ) or type
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		return // no params
	}
	//if here at least one param
	// Loop to handle parameters:
	for {
		//compile the type
		ce.CompileType()
		varType := ce.currentToken.Token_content
		ce.GetToken() //varName
		// Write the variable name.
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name param list " + ce.currentToken.Token_content)
		}
		varName := ce.currentToken.Token_content
		ce.symbolTable.Define(varName, varType, "argument")
		ce.GetToken() //, or )
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ",") {
			break //no more params, returning with )
		}
		ce.GetToken() //if here was comma, advance
	}
}

// CompileSubroutineBody compiles the body of a method, function, or constructor.
func (ce *CompilationEngine) CompileSubroutineBody() {
	// '{' varDec* statements '}'
	ce.GetToken() //open brace
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected { start subroutbody " + ce.currentToken.Token_content)
	}
	ce.GetToken() //get var if var
	// Loop to handle variable declarations (var):
	for {
		if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content != "var" {
			break //no vars!
		}
		//otherwise (go into varDec with Var)
		ce.CompileVarDec() //do we need gettoken after this or in vardec for the loop?
		ce.GetToken()      //left vardec with ; need next token to see if var or not
		//if not it'll be let/if/while/do/return/}
	}
	// Handle bootstrap code for dif subroutine types
	ce.vmWriter.WriteFunction(ce.currentClassName+"."+ce.currentSubroutineName, ce.symbolTable.VarCount("local"))
	if ce.symbolTable.SubroutineKind == "constructor" {
		ce.vmWriter.WritePush("constant", ce.symbolTable.VarCount("field"))
		ce.vmWriter.WriteCall("Memory.alloc", 1)
		ce.vmWriter.WritePop("pointer", 0)
	} else if ce.symbolTable.SubroutineKind == "method" {
		ce.vmWriter.WritePush("argument", 0)
		ce.vmWriter.WritePop("pointer", 0)
	}
	//if no var we have first token of a statement or } end of subroutine
	//otherwise? ; or ?
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
	//came in with var
	// Advance the tokeniser
	ce.GetToken() //type
	//compile the type
	ce.CompileType()
	varType := ce.currentToken.Token_content

	for {
		//compile the type
		ce.GetToken() //varName
		// Write the variable name.
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name vardec loop " + ce.currentToken.Token_content)
		}
		varName := ce.currentToken.Token_content
		ce.symbolTable.Define(varName, varType, "local") //was argument - why??
		ce.GetToken()                                    //, or ;
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ",") {
			break //no more vars, leaving with ;
		}
	}
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected symbol ; end of vardec " + ce.currentToken.Token_content)
	}
}

// CompileStatements compiles a sequence of statements.
func (ce *CompilationEngine) CompileStatements() {
	// (letStatement|ifStatement|whileStatement|doStatement|returnStatement)*
	// enter with token: let/if/while/do/return/}
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
		ce.GetToken() //double check if enough
	}
}

// CompileLet compiles a let statement.
func (ce *CompilationEngine) CompileLet() {
	// 'let' identifier ('[' expression ']')? '=' expression ';'
	//enter with let
	ce.GetToken()                                           //varname
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for var name in let " + ce.currentToken.Token_content)
	}
	varName := ce.currentToken.Token_content
	seg := ce.GetSeg(ce.symbolTable.KindOf(varName))
	index := ce.symbolTable.IndexOf(varName)
	// Advance the tokeniser to check for array indexing: If the current token is an opening bracket [,
	// write the bracket and call compileExpression. Write the closing bracket ] and advance the tokeniser.
	ce.GetToken() // [ or =
	isArray := false
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "[" {
		isArray = true //token is [
		ce.vmWriter.WritePush(seg, index)
		ce.GetToken() // expression
		ce.CompileExpression()
		//check if need get token for ] here
		if ce.currentToken.Token_content == "]" {
		} else {
			ce.GetToken()
		} // expect ]
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
			panic("Unexpected token! Expected ] let " + ce.currentToken.Token_content)
		}
		ce.vmWriter.WriteArithmetic("add")
		ce.GetToken() // for = - check if need
	}
	// Write the equals sign =.
	// NOTE: might not need =?

	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "=") {
		panic("Unexpected token! Expected = ")
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken() //expression
	ce.CompileExpression()

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
		ce.vmWriter.WritePop(seg, index)
	}
}

// CompileDo compiles a do statement.
func (ce *CompilationEngine) CompileDo() {
	// 'do' subroutineCall ';'
	//enter with do
	// Advance the tokeniser and call compileSubroutineCall
	ce.CompileSubroutineCall()
	ce.GetToken() //for ;
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ; Do " + ce.currentToken.Token_content)
	}
	//inside subroutineCall?
	ce.vmWriter.WritePop("temp", 0)
}

// CompileReturn compiles a return statement.
func (ce *CompilationEngine) CompileReturn() {
	// 'return' expression? ';'
	//enter with return
	// Advance the tokeniser to check for an expression: If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken() // expression or ;
	if ce.currentToken.Token_content != ";" {
		ce.CompileExpression()
	} else {
		ce.vmWriter.WritePush("constant", 0)
	}
	ce.vmWriter.WriteReturn()
	//check should be for ;
	// Write the semicolon ;. - get token before leaving compileExpression maybe
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ";") {
		panic("Unexpected token! Expected ; Return" + ce.currentToken.Token_content)
	}
}

// CompileWhile compiles a while statement.
func (ce *CompilationEngine) CompileWhile() {
	// 'while' '('expression')' '{' statements '}'
	//enter with while
	labelContinue := ce.NewLabel()
	labelTop := ce.NewLabel()
	ce.vmWriter.WriteLabel(labelTop)
	ce.GetToken() //get (
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected ( open while " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken() //inside expression? - so far seems to from outside
	ce.CompileExpression()
	//again do we come out with )?
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected ) close while " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteArithmetic("not")
	ce.vmWriter.WriteIfGoto(labelContinue)
	ce.GetToken() //expecting {
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected { open body while " + ce.currentToken.Token_content)
	}
	// Call compileStatements.
	ce.GetToken() // token: let/if/while/do/return/}
	ce.CompileStatements()
	//need get token here? for }?
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected } close body while " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteGoTo(labelTop)
	ce.vmWriter.WriteLabel(labelContinue)
}

// CompileIf compiles an if statement, possibly with a trailing else clause.
func (ce *CompilationEngine) CompileIf() {
	// 'if' '('expression')' '{'statements'}' ('else' '{'statements'}')?
	//enter with if
	labelElse := ce.NewLabel()
	labelEnd := ce.NewLabel()
	ce.GetToken() //for (
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(") {
		panic("Unexpected token! Expected ( open if " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and call compileExpression.
	ce.GetToken() //expression
	ce.CompileExpression()
	// get token before leave compileExpression?
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected ) close if " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteArithmetic("not")
	ce.vmWriter.WriteIfGoto(labelElse)
	// Advance the tokeniser and write the opening brace {.
	ce.GetToken() //expecting {
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
		panic("Unexpected token! Expected { open body if " + ce.currentToken.Token_content)
	}
	// Call compileStatements.
	ce.GetToken() //token: let/if/while/do/return/}
	ce.CompileStatements()
	//need get token here?
	// Write the closing brace }.
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
		panic("Unexpected token! Expected } close body if " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteGoTo(labelEnd)
	ce.vmWriter.WriteLabel(labelElse)
	// Advance the tokeniser to check for an else clause: If the current token is else, write the keyword else,
	// the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken() //else or what?
	if ce.currentToken.Token_type == tokeniser.KEYWORD && ce.currentToken.Token_content == "else" {
		ce.GetToken() //expecting {
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "{") {
			panic("Unexpected token! Expected { else open " + ce.currentToken.Token_content)
		}
		ce.GetToken() //token: let/if/while/do/return/}
		ce.CompileStatements()
		//need gettoken here?
		if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "}") {
			panic("Unexpected token! Expected }")
		}
	} else {
		ce.GoBackToken() //to }
	}
	ce.vmWriter.WriteLabel(labelEnd)
}

// CompileSubroutineCall compiles a subroutine call
func (ce *CompilationEngine) CompileSubroutineCall() {
	// identifier'('expressionList')' | identifier'.'identifier'('expressionList')'
	//enter with do
	ce.GetToken() //subroutineName
	// Write the subroutine name
	if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name call " + ce.currentToken.Token_content)
	}
	name := ce.currentToken.Token_content
	// Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken() //for ( or .
	nArgs := 0
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "(" {
		ce.vmWriter.WritePush("pointer", 0)
		name = ce.currentClassName + "." + name
	} else if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "." {
		ce.GetToken()                                           //next identifier
		if ce.currentToken.Token_type != tokeniser.IDENTIFIER { // not an identifier
			panic("Unexpected token type! Expected identifier for class or var name in subroutinecall with . " + ce.currentToken.Token_content)
		}
		objName := ce.currentToken.Token_content
		//subroutineName := ce.currentToken.Token_content
		typ := ce.symbolTable.TypeOf(objName)
		if typ == "int" || typ == "boolean" || typ == "char" || typ == "void" {
			panic("Unexpected token type! Expected non-built in type in subroutine call " + ce.currentToken.Token_content)
		} else if typ == "" { //This or the symbol table type of need to be matched
			name += "." + objName
		} else {
			nArgs = 1
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(objName)), ce.symbolTable.IndexOf(objName))
			name = ce.symbolTable.TypeOf(objName) + "." + name
		}
		ce.GetToken() //from subroutinename to (
	} else {
		panic("Unexpected token! Expected ( or . in subroutinecall " + ce.currentToken.Token_content)
	}
	// Advance the tokeniser and call compileExpressionList.
	//token at this point is (
	nArgs = ce.CompileExpressionList() //+ 1
	// Write the closing parenthesis ). - when leave expressionList did a get token?
	if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
		panic("Unexpected token! Expected ) end of subroutinecall " + ce.currentToken.Token_content)
	}
	ce.vmWriter.WriteCall(name, nArgs)
}

// CompileExpressionList is responsible for compiling a (possibly empty) comma-separated list of expressions. This list is typically found within the argument list of a subroutine call.
func (ce *CompilationEngine) CompileExpressionList() int {
	// (expression (',' expression)*)?
	//enter with (
	ce.GetToken() // or ) or expression
	nArgs := 0
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		return nArgs // no args
	}
	// If there is at least one expression: Call compileExpression to compile the first expression.
	// Loop to handle additional expressions separated by commas
	ce.CompileExpression() // do getToken before break
	nArgs++
	if ce.currentToken.Token_content == "," {
		for ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
			ce.GetToken()
			ce.CompileExpression() // did a getToken before breaking
			nArgs++
		}
	}
	if ce.currentToken.Token_content == ")" {
	} else {
		ce.GetToken()
	} // comma if but can be )
	// while have comma, another expression
	if ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")" {
		return nArgs // no args
	}
	for ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "," {
		ce.GetToken()
		ce.CompileExpression() // did a getToken before breaking
		nArgs++
		ce.GetToken() // for comma if not from compile expression or something else - end of expressionlist should be a )
	}
	return nArgs
}

// CompileExpression compiles an expression.
func (ce *CompilationEngine) CompileExpression() {
	// term (op term)*
	// Loop to handle additional terms connected by operators:
	for {
		ce.CompileTerm() //what do we leave compile term with?
		ce.GetToken()
		fmt.Println("returnn from compile term to compile expression - cur token is: " + ce.currentToken.Token_content)
		if !(isOperator(ce.currentToken.Token_content)) {
			break
		}
		op := ce.currentToken.Token_content
		ce.GetToken()
		ce.CompileTerm()
		ce.WriteArithmeticCommand(op) //wait until term is placed in stack before putting in op for postfix notation
		ce.GetToken()
		if !(isOperator(ce.currentToken.Token_content)) {
			break
		}
	}
}

// CompileTerm compiles a term.
func (ce *CompilationEngine) CompileTerm() {
	// integerConstant|stringConstant|keywordConstant|identifier|identifier'['expression']'|subroutineCall|
	// '(' expression ')'|unaryOp term

	curtokentyp := ce.currentToken.Token_type
	curtokencontent := ce.currentToken.Token_content

	if curtokentyp == tokeniser.INT_CONST {
		fmt.Println("in int const")
		int_token, err := strconv.Atoi(curtokencontent)
		if err != nil {
			return
		}
		ce.vmWriter.WritePush("constant", int_token)
	} else if curtokentyp == tokeniser.STRING_CONST {
		fmt.Println("in string const")
		stringVal := curtokencontent
		ce.vmWriter.WritePush("constant", len(stringVal))
		ce.vmWriter.WriteCall("String.new", 1) //1 param, length of string
		for _, char := range stringVal {       //NOTE: needs to be rune here?
			ce.vmWriter.WritePush("constant", int(char))
			ce.vmWriter.WriteCall("String.appendChar", 2) //above push and result of new string will be in stack
		}
	} else if curtokentyp == tokeniser.KEYWORD && (curtokencontent == "true" || curtokencontent == "false" || curtokencontent == "null" || curtokencontent == "this") { // got token before called compileTerm
		fmt.Println("in keyword")
		switch curtokencontent {
		case "true":
			ce.vmWriter.WritePush("constant", 0)
			ce.vmWriter.WriteArithmetic("not")
		case "false", "null":
			ce.vmWriter.WritePush("constant", 0)
		case "this":
			ce.vmWriter.WritePush("pointer", 0)
		}
	} else if curtokentyp == tokeniser.IDENTIFIER {
		ce.GetToken()
		nexttokentyp := ce.currentToken.Token_type
		nexttokencontent := ce.currentToken.Token_content
		// save current identifier, and then get next token to see if symbol - look ahead
		if nexttokentyp == tokeniser.SYMBOL && nexttokencontent == "[" {
			// array
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(curtokencontent)), ce.symbolTable.IndexOf(curtokencontent))
			ce.GetToken()          //expression
			ce.CompileExpression() // what do we leave expression with?
			if ce.currentToken.Token_content == "]" {

			} else {
				ce.GetToken() // expect ]
			}
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == "]") {
				panic("Unexpected token! Expected ] in [ of term " + ce.currentToken.Token_content)
			}
			ce.vmWriter.WriteArithmetic("add")
			ce.vmWriter.WritePop("pointer", 1)
			ce.vmWriter.WritePush("that", 0)
		} else if nexttokentyp == tokeniser.SYMBOL && (nexttokencontent == "(" || nexttokencontent == ".") {
			// subroutine call
			// need to go back so function is at subroutine name
			if nexttokencontent == "." {
				ce.GoBackToken() //token now matches curtoken instead of nexttoken
				ce.GoBackToken() //because we go forward a token in subroutinecall
			}
			ce.CompileSubroutineCall()
		} else {
			ce.vmWriter.WritePush(ce.GetSeg(ce.symbolTable.KindOf(curtokencontent)), ce.symbolTable.IndexOf(curtokencontent))
			// need to move token back one since not using current token here
			ce.GoBackToken() //token now matches curtoken instead of nexttoken
		}
	} else if curtokentyp == tokeniser.SYMBOL {
		symbol := curtokencontent
		if symbol == "(" {
			ce.GetToken()          //so not to stay on ( indefinitely
			ce.CompileExpression() //what do we leave compile expression with?
			if !(ce.currentToken.Token_type == tokeniser.SYMBOL && ce.currentToken.Token_content == ")") {
				panic("Unexpected token! Expected ) in compile term symbol " + ce.currentToken.Token_content)
			}
		} else if symbol == "-" || symbol == "~" {
			unaryOp := curtokencontent
			ce.GetToken()
			ce.CompileTerm()
			if unaryOp == "-" {
				ce.vmWriter.WriteArithmetic("neg")
			} else if unaryOp == "~" {
				ce.vmWriter.WriteArithmetic("not")
			}
		}
	}
}

func (ce *CompilationEngine) CompileType() {
	if !(ce.currentToken.Token_type == tokeniser.KEYWORD && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) && !(ce.currentToken.Token_type == tokeniser.IDENTIFIER) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for type or identifier")
	}
}

func (ce *CompilationEngine) GetToken() {
	ce.currentToken = &ce.tokeniser.Tokens[ce.currentTokenIndex]
	ce.currentTokenIndex = ce.currentTokenIndex + 1
	//PrintCaller()
	//fmt.Println(strconv.Itoa(ce.currentTokenIndex) + " " + ce.currentToken.Token_content)
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

func (ce *CompilationEngine) NewLabel() string {
	var l = ce.labelIndex
	ce.labelIndex = l + 1
	return "LABEL_" + strconv.Itoa(l)
}

func (ce *CompilationEngine) GoBackToken() {
	ce.currentTokenIndex = ce.currentTokenIndex - 1
	ce.currentToken = &ce.tokeniser.Tokens[ce.currentTokenIndex-1]
	//PrintCaller()
	//fmt.Println(strconv.Itoa(ce.currentTokenIndex) + " " + ce.currentToken.Token_content)
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

func isOperator(symbol string) bool {
	return symbol == "+" || symbol == "-" || symbol == "*" || symbol == "/" ||
		symbol == "&" || symbol == "|" || symbol == "<" || symbol == ">" || symbol == "="
}

func PrintCaller() {
	_, _, line, ok := runtime.Caller(2)
	if !ok {
		fmt.Println("Could not get caller information")
		return
	}
	fmt.Printf("Called from line %d\n", line)
}
