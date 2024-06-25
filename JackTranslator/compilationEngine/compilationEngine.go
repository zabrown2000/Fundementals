package compilationEngine

import (
	"Fundementals/JackTranslator/tokeniser"
	"bufio"
	"os"
)

// new method: tokeniser creates list and sends it here where it gets handled
// will have function to get next token in tokens list
// need 2 writings:
// 1. xml with just list of keywords, symbols, identifiers, etc w/o enclosing tags
// 2. xml with hierarchy

type CompilationEngine struct {
	tokeniser         *tokeniser.Tokeniser
	plainWriter       *bufio.Writer
	hierarchWriter    *bufio.Writer
	currentToken      *tokeniser.Token
	currentTokenIndex int
}

func NewCompilationEngine(plainOutputFile string, hierarchOutputFile string, tokeniser *tokeniser.Tokeniser) *CompilationEngine {
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

func (ce *CompilationEngine) CompileClass() {

	//Purpose: Compiles a complete class.
	//Steps:
	//1. Write the opening tag <class>.
	// TC and close tag for plain when nothing's been opened?
	ce.WriteOpenTag(ce.hierarchWriter, "class")
	ce.WriteCloseTag(ce.plainWriter, "tokens")
	//2. Advance the tokeniser to the next token and expect the keyword class.
	ce.GetToken()
	//TC changing if to compare content of token not just type
	//if ce.currentToken.Token_type != 1 { // not a keyword
	if !(ce.currentToken.Token_type == 1 && ce.currentToken.Token_content == "class") { // not a keyword
		panic("Unexpected token type! Expected keyword")
	}
	//3. Write the class keyword.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//4. Advance the tokeniser and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier")
	}
	//5. Write the class name.
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//6. Advance the tokeniser and expect the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token_content != "{" { // not an identifier
		panic("Unexpected token! Expected {")
	}
	//7. Write the { symbol.
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//8. Loop to handle class variable declarations (static or field) and subroutine declarations (constructor, function, or method):
	//     If the current token is static or field, call compileClassVarDec.
	//     If the current token is constructor, function, or method, call compileSubroutine.
	//     Otherwise, break the loop.
	for { // TC before now, you had getToken before all comparisons of content - this time not - intentional? if not, outside or in the for loop?
		ce.GetToken()
		if ce.currentToken.Token_type == 1 && (ce.currentToken.Token_content == "static" || ce.currentToken.Token_content == "field") {
			// TC this is for inside CompileClassVarDec - but you didn't write it to static or field to file before advancing
			//so when you go to write you lost it - moving to inside CompileClassVarDec
			//ce.GetToken()
			ce.CompileClassVarDec()
		} else if ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method" {
			ce.CompileSubroutine()
		} else {
			break // class is complete
		}
	}
	//TC also need to check for the closing brace - with get token can't just assume it's there
	//9. Write the closing brace } symbol.
	if ce.currentToken.Token_content != "}" { // not an identifier
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", "}")
	ce.WriteXML(ce.plainWriter, "symbol", "}")
	//10. Write the closing tag </class>.
	ce.WriteCloseTag(ce.hierarchWriter, "class")
	ce.WriteCloseTag(ce.plainWriter, "tokens")
}

func (ce *CompilationEngine) CompileClassVarDec() {

	//Purpose: Compiles a static declaration or a field declaration.
	//Steps:
	//1. Write the opening tag <classVarDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "classVarDec")
	//2. Write the current token (static or field).
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content) //no need to check if field or static, did that in compileClass
	//3. Advance the tokeniser and write the type (e.g., int, boolean, or a class name).
	// TC didn't actually get the next token here
	ce.GetToken()
	// TC also class name is an identifier not a keyword, and you still need to make sure if its specifically int/boolean/char if it's a keyword
	//ce.tokeniser.Advance()
	// TC technically we should be calling CompileType and CompileClassName not just checking here
	if !(ce.currentToken.Token_type == 1 && (ce.currentToken.Token_content == "int" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == 3) { // not a keyword or an identifier
		panic("Unexpected token type! Expected keyword for var type or identifier")
	}
	if ce.currentToken.Token_type == 1 {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == 3 {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	}
	//4. Advance the tokeniser and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	//6. Write the semicolon ; symbol.
	// TC need to check if we got the semicolon not just assume!
	if ce.currentToken.Token_content != ";" {
		panic("Unexpected token type! Expected symbol ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ";")
	ce.WriteXML(ce.plainWriter, "symbol", ";")
	//7. Write the closing tag </classVarDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "classVarDec")
}

func (ce *CompilationEngine) CompileSubroutine() {
	//Purpose: Compiles a complete method, function, or constructor.
	//Steps:
	//1. Write the opening tag <subroutineDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "subroutineDec")
	//2. Write the current token (constructor, function, or method).
	if ce.currentToken.Token_type == 1 && (ce.currentToken.Token_content == "constructor" || ce.currentToken.Token_content == "function" || ce.currentToken.Token_content == "method") {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	//3. Advance the tokeniser and write the return type (void or a type).
	// TC again we should technically be calling CompileType (reminder types are either 3 specific keywords or an identifier
	ce.GetToken()
	if !(ce.currentToken.Token_type == 1 && (ce.currentToken.Token_content == "void" || ce.currentToken.Token_content == "char" ||
		ce.currentToken.Token_content == "boolean")) || !(ce.currentToken.Token_type == 3) { // not a keyword or a type
		panic("Unexpected token type! Expected keyword or identifier for subroutine return type")
	}
	if ce.currentToken.Token_type == 1 {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	} else if ce.currentToken.Token_type == 3 {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	}
	//4. Advance the tokeniser and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type == 3 {
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected identifier for subroutine name")
	}

	//5. Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "(" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token! Expected (")
	}
	//6. Advance the tokeniser and call compileParameterList.
	// TC need to check if next token is ) not ( as not all functions/methods/constructors actually have a parameter list?
	ce.GetToken()
	if ce.currentToken.Token_content != ")" {
		ce.CompileParameterList()
	}
	//7. Write the closing parenthesis ) - when no parameters token from getToken above, otherwise from getToken in broken loop in CompileParameterList
	if ce.currentToken.Token_content == ")" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	} else {
		panic("Unexpected token type! Expected )")
	}
	//8. Advance the tokeniser and call compileSubroutineBody.
	// TC for sake of uniformity - call getToken from inside CompileSubroutineBody? - already called in CompileParameterList -
	// TC this will advance twice before we check anything so removing here
	//ce.GetToken()
	ce.CompileSubroutineBody()
	//9. Write the closing tag </subroutineDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineDec")
}

func (ce *CompilationEngine) CompileParameterList() {
	//Purpose: Compiles a parameter list.
	//Steps:
	//1. Write the opening tag <parameterList>
	ce.WriteOpenTag(ce.hierarchWriter, "parameterList")
	//2. Loop to handle parameters:
	//     Write the type and variable name for each parameter.
	//     If the current token is a comma ,, write the comma and process the next parameter.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type != 1 { // not a keyword
			panic("Unexpected token type! Expected keyword for var type")
		}
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		//4. Advance the tokeniser and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
		ce.GetToken()
		if ce.currentToken.Token_type != 3 { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name")
		}
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
		} else {
			break //no more params
		}
	}
	//3. Write the closing tag </parameterList>.
	ce.WriteCloseTag(ce.hierarchWriter, "parameterList")
}

func (ce *CompilationEngine) CompileSubroutineBody() {
	//Purpose: Compiles the body of a method, function, or constructor.
	//Steps:
	//1. Write the opening tag <subroutineBody>.
	ce.WriteOpenTag(ce.hierarchWriter, "subroutineBody")
	//2. Write the opening brace {.
	if ce.currentToken.Token_content != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//3. Loop to handle variable declarations (var):
	//     If the current token is var, call compileVarDec.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_content == "var" {
			ce.CompileVarDec()
		} else {
			break //no more vars
		}
	}
	//4. Call compileStatements to handle the statements within the subroutine body.
	ce.GetToken()
	ce.CompileStatements()
	//5. Write the closing brace }.
	if ce.currentToken.Token_content != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//6. Write the closing tag </subroutineBody>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineBody")
}

func (ce *CompilationEngine) CompileVarDec() {
	//Purpose: Compiles a var declaration.
	//Steps:
	//1. Write the opening tag <varDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "varDec")
	//2. Write the current token var.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser and write the type.
	if ce.currentToken.Token_type != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for var type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//4. Advance the tokeniser and write the first variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
			if ce.currentToken.Token_type != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
		} else {
			break // no more variables
		}
	}
	//6. Write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token_content != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//7. Write the closing tag </varDec>.
	ce.WriteCloseTag(ce.hierarchWriter, "varDec")
}

func (ce *CompilationEngine) CompileStatements() {
	//Purpose: Compiles a sequence of statements.
	//Steps:
	//1. Write the opening tag <statements>.
	ce.WriteOpenTag(ce.hierarchWriter, "statements")
	//2. Loop to handle each statement:
	//      Based on the current token, call the appropriate function (compileLet, compileIf, compileWhile, compileDo, or compileReturn).
	//      If none of these keywords match, break the loop.
	for {
		if ce.currentToken.Token_content == "let" {
			ce.GetToken()
			ce.CompileLet()
		} else if ce.currentToken.Token_content == "if" {
			ce.GetToken()
			ce.CompileIf()
		} else if ce.currentToken.Token_content == "while" {
			ce.GetToken()
			ce.CompileWhile()
		} else if ce.currentToken.Token_content == "do" {
			ce.GetToken()
			ce.CompileDo()
		} else if ce.currentToken.Token_content == "return" {
			ce.GetToken()
			ce.CompileReturn()
		} else {
			break
		}
	}
	//3. Write the closing tag </statements>
	ce.WriteCloseTag(ce.hierarchWriter, "statements")

}

func (ce *CompilationEngine) CompileDo() {
	//Purpose: Compiles a do statement.
	//Steps:
	//1. Write the opening tag <doStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "doStatement")
	//2. Write the current token do.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//4. Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_content != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//5. Advance the tokeniser and call compileExpressionList.
	ce.GetToken()
	ce.CompileExpressionList()
	//6. Write the closing parenthesis ).
	ce.GetToken() //----------need this?
	if ce.currentToken.Token_content != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//7. Advance the tokeniser and write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token_content != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//8. Write the closing tag </doStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "doStatement")
}

func (ce *CompilationEngine) CompileLet() {
	//Purpose: Compiles a let statement.
	//Steps:
	//1. Write the opening tag <letStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "letStatement")
	//2. Write the current token let.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser and write the variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for var name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token_content)
	//4. Advance the tokeniser to check for array indexing:
	//      If the current token is an opening bracket [, write the bracket and call compileExpression.
	//      Write the closing bracket ] and advance the tokeniser.
	ce.GetToken()
	if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "[" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.GetToken()
		ce.CompileExpression()
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "]" {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
			ce.GetToken()
		}
	}
	//5. Write the equals sign =.
	if ce.currentToken.Token_content != "=" {
		panic("Unexpected token! Expected =")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//6. Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//7. Write the semicolon ;.
	if ce.currentToken.Token_content != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//8. Write the closing tag </letStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "letStatement")
}

func (ce *CompilationEngine) CompileWhile() {
	//Purpose: Compiles a while statement.
	//Steps:
	//1. Write the opening tag <whileStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "whileStatement")
	//2. Write the current token while.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_content != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//4. Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//5. Write the closing parenthesis ).
	ce.GetToken()
	if ce.currentToken.Token_content != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//6. Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token_content != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//7. Advance the tokeniser and call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	//8. Write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_content != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//9. Write the closing tag </whileStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "whileStatement")
}

func (ce *CompilationEngine) CompileReturn() {
	//Purpose: Compiles a return statement.
	//Steps:
	//1. Write the opening tag <returnStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "returnStatement")
	//2. Write the current token return.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser to check for an expression:
	//     If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken()
	if ce.currentToken.Token_content != ";" {
		ce.CompileExpression()
	}
	//4. Write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token_content != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//5. Write the closing tag </returnStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "returnStatement")
}

func (ce *CompilationEngine) CompileIf() {
	//Purpose: Compiles an if statement, possibly with a trailing else clause.
	//Steps:
	//1. Write the opening tag <ifStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "ifStatement")
	//2. Write the current token if.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
	//3. Advance the tokeniser and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_content != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//4. Advance the tokeniser and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//5. Write the closing parenthesis ).
	ce.GetToken()
	if ce.currentToken.Token_content != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//6. Advance the tokeniser and write the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token_content != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//7. Advance the tokeniser and call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	//8. Write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_content != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	//9. Advance the tokeniser to check for an else clause:
	//      If the current token is else, write the keyword else, the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token_content == "else" {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		ce.GetToken()
		if ce.currentToken.Token_content != "{" {
			panic("Unexpected token! Expected {")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
		ce.GetToken()
		ce.CompileStatements()
		ce.GetToken()
		if ce.currentToken.Token_content != "}" {
			panic("Unexpected token! Expected }")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content)
	}
	//10. Write the closing tag </ifStatement>
	ce.WriteCloseTag(ce.hierarchWriter, "ifStatement")
}

func (ce *CompilationEngine) CompileExpression() {
	//Purpose: Compiles an expression.
	//Steps:
	//1. Write the opening tag <expression>.
	ce.WriteOpenTag(ce.hierarchWriter, "expression")
	//2. Call compileTerm. --advance before calling callexpression
	ce.CompileTerm()
	//3. Loop to handle additional terms connected by operators:
	//      If the current token is an operator, write the operator and call compileTerm for the next term.
	//      Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && (ce.currentToken.Token_content == "+" || ce.currentToken.Token_content == "-" || ce.currentToken.Token_content == "*" || ce.currentToken.Token_content == "/" || ce.currentToken.Token_content == "&" || ce.currentToken.Token_content == "|" || ce.currentToken.Token_content == "<" || ce.currentToken.Token_content == ">" || ce.currentToken.Token_content == "=") {
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
	//4. Write the closing tag </expression>
	ce.WriteCloseTag(ce.hierarchWriter, "expression")
}

func (ce *CompilationEngine) CompileTerm() {
	//Purpose: Compiles a term.
	//Steps:
	//1. Write the opening tag <term>.
	ce.WriteOpenTag(ce.hierarchWriter, "term")
	//2. Depending on the current token, handle different types of terms:
	//      If the token is an integer constant, write the integer constant.
	//      If the token is a string constant, write the string constant.
	//      If the token is a keyword constant, write the keyword.
	//      If the token is an identifier, handle variable names, array entries, or subroutine calls.
	//      If the token is an opening parenthesis (, call compileExpression and write the closing parenthesis ).
	//      If the token is a unary operator, write the operator and call compileTerm.
	if ce.currentToken.Token_type == 1 {
		// keyword
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token_content)
		ce.GetToken()
	} else if ce.currentToken.Token_type == 2 {
		// identifier
		// save current id, and then get next to see if symbol - look ahead
		identifier := ce.currentToken.Token_content
		ce.GetToken()

		if ce.currentToken.Token_type == 2 { //symbol
			symbol := ce.currentToken.Token_content
			if symbol == "[" {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // arrayName
				ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
				ce.WriteXML(ce.plainWriter, "symbol", symbol) // [
				ce.GetToken()
				ce.CompileExpression() // end with break after getToken
				if ce.currentToken.Token_content != "]" {
					panic("Unexpected token! Expected ]")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // [
				ce.GetToken()
			} else if symbol == "(" {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // subroutineName
				ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
				ce.WriteXML(ce.plainWriter, "symbol", symbol) // (
				ce.GetToken()
				ce.CompileExpressionList() // end with break after getToken
				if ce.currentToken.Token_content != ")" {
					panic("Unexpected token! Expected )")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // )
				ce.GetToken()
			} else if symbol == "." {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // className or varName
				ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
				ce.WriteXML(ce.plainWriter, "symbol", symbol) // .
				ce.GetToken()
				if ce.currentToken.Token_type != 3 {
					panic("Unexpected token! Expected identifier for subroutine name ")
				}
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // subroutineName
				ce.GetToken()
				if ce.currentToken.Token_content != "(" {
					panic("Unexpected token! Expected (")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // (
				ce.GetToken()
				ce.CompileExpressionList()
				if ce.currentToken.Token_content != ")" {
					panic("Unexpected token! Expected )")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // )
				ce.GetToken()
			} else {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier)
				// current token is the one after the identifier
				// NOTE: need function to move token back?
			}
		} else {
			ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
			ce.WriteXML(ce.plainWriter, "identifier", identifier)
		}
		// add ce.GetToken() here?
	} else if ce.currentToken.Token_type == 3 {
		// symbol
		symbol := ce.currentToken.Token_content
		if symbol == "(" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileExpression()
			if ce.currentToken.Token_content != ")" {
				panic("Unexpected token! Expected )")
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // )
			ce.GetToken()
		} else if symbol == "-" || symbol == "~" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileTerm()
		}
	} else if ce.currentToken.Token_type == 4 {
		// int constant
		ce.WriteXML(ce.hierarchWriter, "int_const", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "int_const", ce.currentToken.Token_content)
		ce.GetToken()
	} else if ce.currentToken.Token_type == 5 {
		// string constant
		ce.WriteXML(ce.hierarchWriter, "string_const", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "string_const", ce.currentToken.Token_content)
		ce.GetToken()
	}
	//3. Write the closing tag </term>.
	ce.WriteCloseTag(ce.hierarchWriter, "term")
}

func (ce *CompilationEngine) CompileExpressionList() {
	//Purpose: The compileExpressionList function is responsible for compiling a (possibly empty) comma-separated list of expressions. This list is typically found within the argument list of a subroutine call.
	//Steps:
	//1. Write the opening tag <expressionList>.
	ce.WriteOpenTag(ce.hierarchWriter, "expressionList")
	//2. Check if the current token indicates the start of an expression. This can be identified by looking for tokens that can start an expression such as integer constants, string constants, keyword constants, variable names, subroutine calls, expressions enclosed in parentheses, and unary operators.
	if ce.currentToken.Token_content == ")" {
		return // empty list
	}
	//3. If there is at least one expression:
	//     Call compileExpression to compile the first expression.
	//     Loop to handle additional expressions separated by commas:
	//          If the current token is a comma ,, write the comma symbol.
	//          Advance the tokeniser.
	//          Call compileExpression to compile the next expression.
	// Check if there is at least one expression to compile
	ce.CompileExpression() // do getToken before break

	// while have comma, another expression
	for ce.currentToken.Token_type == 2 && ce.currentToken.Token_content == "," {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token_content)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token_content) // ','
		ce.GetToken()
		ce.CompileExpression() // did a getToken before breaking
	}
	//4. Write the closing tag </expressionList>.
	ce.WriteCloseTag(ce.hierarchWriter, "expressionList")
}

// WriteXML helper functions
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
