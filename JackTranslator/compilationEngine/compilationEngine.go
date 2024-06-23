package compilationEngine

import (
	"Fundementals/JackTranslator/tokeniser"
	"bufio"
	"os"
)

// new method: tokenizer creates list and sends it here where it gets handled
// will have function to get next token in tokens list
// need 2 writings:
// 1. xml with just list of keywords, symbols, indentifiers, etc w/o enclosing tags
// 2. xml with hierarchy

type CompilationEngine struct {
	tokenizer         *tokeniser.Tokeniser
	plainWriter       *bufio.Writer
	hierarchWriter    *bufio.Writer
	currentToken      *tokeniser.Token
	currentTokenIndex int
}

func NewCompilationEngine(plainOutputFile string, hierarchOutputFile string, tokenizer *tokeniser.Tokeniser) *CompilationEngine {
	plainFile, err := os.Create(plainOutputFile)
	if err != nil {
		return nil
	}
	plainWriter := bufio.NewWriter(plainFile)

	hierarchFile, err := os.Create(hierarchOutputFile)
	if err != nil {
		return nil
	}
	hierarchWriter := bufio.NewWriter(hierarchFile)

	return &CompilationEngine{
		tokenizer:         tokenizer,
		plainWriter:       plainWriter,
		hierarchWriter:    hierarchWriter,
		currentTokenIndex: 0,
	}
}

func (ce *CompilationEngine) CompileClass() {

	//Purpose: Compiles a complete class.
	//Steps:
	//1. Write the opening tag <class>.
	ce.WriteOpenTag(ce.hierarchWriter, "class")
	ce.WriteCloseTag(ce.plainWriter, "tokens")
	//2. Advance the tokenizer to the next token and expect the keyword class.
	ce.GetToken()
	if ce.currentToken.Token_type != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword")
	}
	//3. Write the class keyword.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//4. Advance the tokenizer and expect the class name (identifier).
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier")
	}
	//5. Write the class name.
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//6. Advance the tokenizer and expect the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token != "{" { // not an identifier
		panic("Unexpected token! Expected {")
	}
	//7. Write the { symbol.
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//8. Loop to handle class variable declarations (static or field) and subroutine declarations (constructor, function, or method):
	//     If the current token is static or field, call compileClassVarDec.
	//     If the current token is constructor, function, or method, call compileSubroutine.
	//     Otherwise, break the loop.
	for ce.tokenizer.HasMoreTokens() {
		if ce.currentToken.Token == "static" || ce.currentToken.Token == "field" {
			ce.GetToken()
			ce.CompileClassVarDec()
		} else if ce.currentToken.Token == "constructor" || ce.currentToken.Token == "function" || ce.currentToken.Token == "method" {
			ce.CompileSubroutine()
		} else {
			break // class is complete
		}
	}
	//9. Write the closing brace } symbol.
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
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token) //no need to check if field or static, did that in compileClass
	//3. Advance the tokenizer and write the type (e.g., int, boolean, or a class name).
	//ce.tokenizer.Advance()
	if ce.currentToken.Token_type != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for var type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//4. Advance the tokenizer and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
			ce.GetToken()
			if ce.currentToken.Token_type != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
		} else {
			break // no more variables
		}
	}
	//6. Write the semicolon ; symbol.
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
	if ce.currentToken.Token_type == 1 && (ce.currentToken.Token == "constructor" || ce.currentToken.Token == "function" || ce.currentToken.Token == "method") {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	} else {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	//3. Advance the tokenizer and write the return type (void or a type).
	ce.GetToken()
	if ce.currentToken.Token_type != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for subroutine return type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//4. Advance the tokenizer and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//5. Advance the tokenizer and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "(" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	} else {
		panic("Unexpected token! Expected (")
	}
	//6. Advance the tokenizer and call compileParameterList.
	ce.GetToken()
	if ce.currentToken.Token == "(" {
		ce.CompileParameterList()
	}
	//7. Write the closing parenthesis ).
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//8. Advance the tokenizer and call compileSubroutineBody.
	ce.GetToken()
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
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
		//4. Advance the tokenizer and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
		ce.GetToken()
		if ce.currentToken.Token_type != 3 { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name")
		}
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
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
	if ce.currentToken.Token != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//3. Loop to handle variable declarations (var):
	//     If the current token is var, call compileVarDec.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token == "var" {
			ce.CompileVarDec()
		} else {
			break //no more vars
		}
	}
	//4. Call compileStatements to handle the statements within the subroutine body.
	ce.GetToken()
	ce.CompileStatements()
	//5. Write the closing brace }.
	if ce.currentToken.Token != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//6. Write the closing tag </subroutineBody>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineBody")
}

func (ce *CompilationEngine) CompileVarDec() {
	//Purpose: Compiles a var declaration.
	//Steps:
	//1. Write the opening tag <varDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "varDec")
	//2. Write the current token var.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer and write the type.
	if ce.currentToken.Token_type != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for var type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//4. Advance the tokenizer and write the first variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
			ce.GetToken()
			if ce.currentToken.Token_type != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
		} else {
			break // no more variables
		}
	}
	//6. Write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
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
		if ce.currentToken.Token == "let" {
			ce.GetToken()
			ce.CompileLet()
		} else if ce.currentToken.Token == "if" {
			ce.GetToken()
			ce.CompileIf()
		} else if ce.currentToken.Token == "while" {
			ce.GetToken()
			ce.CompileWhile()
		} else if ce.currentToken.Token == "do" {
			ce.GetToken()
			ce.CompileDo()
		} else if ce.currentToken.Token == "return" {
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
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer and write the subroutine name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//4. Advance the tokenizer and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//5. Advance the tokenizer and call compileExpressionList.
	ce.GetToken()
	ce.CompileExpressionList()
	//6. Write the closing parenthesis ).
	ce.GetToken() //----------need this?
	if ce.currentToken.Token != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//7. Advance the tokenizer and write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//8. Write the closing tag </doStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "doStatement")
}

func (ce *CompilationEngine) CompileLet() {
	//Purpose: Compiles a let statement.
	//Steps:
	//1. Write the opening tag <letStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "letStatement")
	//2. Write the current token let.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer and write the variable name.
	ce.GetToken()
	if ce.currentToken.Token_type != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for var name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "identifier", ce.currentToken.Token)
	//4. Advance the tokenizer to check for array indexing:
	//      If the current token is an opening bracket [, write the bracket and call compileExpression.
	//      Write the closing bracket ] and advance the tokenizer.
	ce.GetToken()
	if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "[" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
		ce.GetToken()
		ce.CompileExpression()
		ce.GetToken()
		if ce.currentToken.Token_type == 2 && ce.currentToken.Token == "]" {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
			ce.GetToken()
		}
	}
	//5. Write the equals sign =.
	if ce.currentToken.Token != "=" {
		panic("Unexpected token! Expected =")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//6. Advance the tokenizer and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//7. Write the semicolon ;.
	if ce.currentToken.Token != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//8. Write the closing tag </letStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "letStatement")
}

func (ce *CompilationEngine) CompileWhile() {
	//Purpose: Compiles a while statement.
	//Steps:
	//1. Write the opening tag <whileStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "whileStatement")
	//2. Write the current token while.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//4. Advance the tokenizer and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//5. Write the closing parenthesis ).
	ce.GetToken()
	if ce.currentToken.Token != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//6. Advance the tokenizer and write the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//7. Advance the tokenizer and call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	//8. Write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//9. Write the closing tag </whileStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "whileStatement")
}

func (ce *CompilationEngine) CompileReturn() {
	//Purpose: Compiles a return statement.
	//Steps:
	//1. Write the opening tag <returnStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "returnStatement")
	//2. Write the current token return.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer to check for an expression:
	//     If the current token is not a semicolon ;, call compileExpression.
	ce.GetToken()
	if ce.currentToken.Token != ";" {
		ce.CompileExpression()
	}
	//4. Write the semicolon ;.
	ce.GetToken()
	if ce.currentToken.Token != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//5. Write the closing tag </returnStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "returnStatement")
}

func (ce *CompilationEngine) CompileIf() {
	//Purpose: Compiles an if statement, possibly with a trailing else clause.
	//Steps:
	//1. Write the opening tag <ifStatement>.
	ce.WriteOpenTag(ce.hierarchWriter, "ifStatement")
	//2. Write the current token if.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
	//3. Advance the tokenizer and write the opening parenthesis (.
	ce.GetToken()
	if ce.currentToken.Token != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//4. Advance the tokenizer and call compileExpression.
	ce.GetToken()
	ce.CompileExpression()
	//5. Write the closing parenthesis ).
	ce.GetToken()
	if ce.currentToken.Token != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//6. Advance the tokenizer and write the opening brace {.
	ce.GetToken()
	if ce.currentToken.Token != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//7. Advance the tokenizer and call compileStatements.
	ce.GetToken()
	ce.CompileStatements()
	//8. Write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
	ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
	//9. Advance the tokenizer to check for an else clause:
	//      If the current token is else, write the keyword else, the opening brace {, call compileStatements, and write the closing brace }.
	ce.GetToken()
	if ce.currentToken.Token == "else" {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
		ce.GetToken()
		if ce.currentToken.Token != "{" {
			panic("Unexpected token! Expected {")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
		ce.GetToken()
		ce.CompileStatements()
		ce.GetToken()
		if ce.currentToken.Token != "}" {
			panic("Unexpected token! Expected }")
		}
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token)
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
		if ce.currentToken.Token_type == 2 && (ce.currentToken.Token == "+" || ce.currentToken.Token == "-" || ce.currentToken.Token == "*" || ce.currentToken.Token == "/" || ce.currentToken.Token == "&" || ce.currentToken.Token == "|" || ce.currentToken.Token == "<" || ce.currentToken.Token == ">" || ce.currentToken.Token == "=") {
			str := ce.currentToken.Token
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
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "keyword", ce.currentToken.Token)
		ce.GetToken()
	} else if ce.currentToken.Token_type == 2 {
		// identifier
		// save current id, and then get next to see if symbol - look ahead
		identifier := ce.currentToken.Token
		ce.GetToken()

		if ce.currentToken.Token_type == 2 { //symbol
			symbol := ce.currentToken.Token
			if symbol == "[" {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // arrayName
				ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
				ce.WriteXML(ce.plainWriter, "symbol", symbol) // [
				ce.GetToken()
				ce.CompileExpression() // end with break after getToken
				if ce.currentToken.Token != "]" {
					panic("Unexpected token! Expected ]")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // [
				ce.GetToken()
			} else if symbol == "(" {
				ce.WriteXML(ce.hierarchWriter, "identifier", identifier)
				ce.WriteXML(ce.plainWriter, "identifier", identifier) // subroutineName
				ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
				ce.WriteXML(ce.plainWriter, "symbol", symbol) // (
				ce.GetToken()
				ce.CompileExpressionList() // end with break after getToken
				if ce.currentToken.Token != ")" {
					panic("Unexpected token! Expected )")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // )
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
				if ce.currentToken.Token != "(" {
					panic("Unexpected token! Expected (")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // (
				ce.GetToken()
				ce.CompileExpressionList()
				if ce.currentToken.Token != ")" {
					panic("Unexpected token! Expected )")
				}
				ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
				ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // )
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
		symbol := ce.currentToken.Token
		if symbol == "(" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileExpression()
			if ce.currentToken.Token != ")" {
				panic("Unexpected token! Expected )")
			}
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
			ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // )
			ce.GetToken()
		} else if symbol == "-" || symbol == "~" {
			ce.WriteXML(ce.hierarchWriter, "symbol", symbol)
			ce.WriteXML(ce.plainWriter, "symbol", symbol)
			ce.GetToken()
			ce.CompileTerm()
		}
	} else if ce.currentToken.Token_type == 4 {
		// int constant
		ce.WriteXML(ce.hierarchWriter, "int_const", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "int_const", ce.currentToken.Token)
		ce.GetToken()
	} else if ce.currentToken.Token_type == 5 {
		// string constant
		ce.WriteXML(ce.hierarchWriter, "string_const", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "string_const", ce.currentToken.Token)
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
	if ce.currentToken.Token == ")" {
		return // empty list
	}
	//3. If there is at least one expression:
	//     Call compileExpression to compile the first expression.
	//     Loop to handle additional expressions separated by commas:
	//          If the current token is a comma ,, write the comma symbol.
	//          Advance the tokenizer.
	//          Call compileExpression to compile the next expression.
	// Check if there is at least one expression to compile
	ce.CompileExpression() // do getToken before break

	// while have comma, another expression
	for ce.currentToken.Token_type == 2 && ce.currentToken.Token == "," {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.currentToken.Token)
		ce.WriteXML(ce.plainWriter, "symbol", ce.currentToken.Token) // ','
		ce.GetToken()
		ce.CompileExpression() // did a getToken before breaking
	}
	//4. Write the closing tag </expressionList>.
	ce.WriteCloseTag(ce.hierarchWriter, "expressionList")
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
	ce.currentToken = &ce.tokenizer.Tokens[ce.currentTokenIndex]
	ce.currentTokenIndex = ce.currentTokenIndex + 1
}
