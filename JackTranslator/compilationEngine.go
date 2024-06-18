package JackTranslator

import (
	"bufio"
	"fmt"
	"os"
)

// new method: tokenizer creates list and sends it here where it gets handled
// will have function to get next token in tokens list
// need 2 writings:
// 1. xml with just list of keywords, symbols, indentifiers, etc w/o enclosing tags
// 2. xml with hierarchy

type CompilationEngine struct {
	tokenizer      *JackTokenizer
	plainWriter    *bufio.Writer
	hierarchWriter *bufio.Writer
	// add current token type
}

// throws error if illegal syntax, based on rules/grammar
// use grammar slides (not lexical elements)

func NewCompilationEngine(plainOutputFile string, hierarchOutputFile string, tokenizer *JackTokenizer) *CompilationEngine {
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
		tokenizer:      tokenizer,
		plainWriter:    plainWriter,
		hierarchWriter: hierarchWriter,
	}
}

func (ce *CompilationEngine) CompileClass() {

	//Purpose: Compiles a complete class.
	//Steps:
	//1. Write the opening tag <class>.
	ce.WriteOpenTag(ce.hierarchWriter, "class")
	ce.WriteCloseTag(ce.plainWriter, "tokens")
	//2. Advance the tokenizer to the next token and expect the keyword class.
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword")
	}
	//3. Write the class keyword.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//4. Advance the tokenizer and expect the class name (identifier).
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier")
	}
	//5. Write the class name.
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
	//6. Advance the tokenizer and expect the opening brace {.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken != "{" { // not an identifier
		panic("Unexpected token! Expected {")
	}
	//7. Write the { symbol.
	ce.WriteXML(ce.hierarchWriter, "symbol", "{")
	ce.WriteXML(ce.plainWriter, "symbol", "{")
	//8. Loop to handle class variable declarations (static or field) and subroutine declarations (constructor, function, or method):
	//     If the current token is static or field, call compileClassVarDec.
	//     If the current token is constructor, function, or method, call compileSubroutine.
	//     Otherwise, break the loop.
	for ce.tokenizer.HasMoreTokens() {
		if ce.tokenizer.currentToken == "static" || ce.tokenizer.currentToken == "field" {
			ce.tokenizer.Advance()
			ce.CompileClassVarDec()
		} else if ce.tokenizer.currentToken == "constructor" || ce.tokenizer.currentToken == "function" || ce.tokenizer.currentToken == "method" {
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
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken) //no need to check if field or static, did that in compileClass
	//3. Advance the tokenizer and write the type (e.g., int, boolean, or a class name).
	//ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for var type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//4. Advance the tokenizer and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.tokenizer.Advance()
		if ce.tokenizer.currentToken == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
			ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
			ce.tokenizer.Advance()
			if ce.tokenizer.TokenType() != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
			ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
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
	if ce.tokenizer.currentToken == "constructor" || ce.tokenizer.currentToken == "function" || ce.tokenizer.currentToken == "method" {
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
		ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	} else {
		panic("Unexpected token type! Expected keyword for subroutine")
	}
	//3. Advance the tokenizer and write the return type (void or a type).
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for subroutine return type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//4. Advance the tokenizer and write the subroutine name.
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
	//5. Advance the tokenizer and write the opening parenthesis (.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken == "(" {
		ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
		ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	} else {
		panic("Unexpected token! Expected (")
	}
	//6. Advance the tokenizer and call compileParameterList.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken == "(" {
		ce.CompileParameterList()
	}
	//7. Write the closing parenthesis ).
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//8. Advance the tokenizer and call compileSubroutineBody.
	ce.tokenizer.Advance()
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
		ce.tokenizer.Advance()
		if ce.tokenizer.TokenType() != 1 { // not a keyword
			panic("Unexpected token type! Expected keyword for var type")
		}
		ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
		ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
		//4. Advance the tokenizer and write the first variable name. ---are all 3 written consec inside tag? seems like do writeXml
		ce.tokenizer.Advance()
		if ce.tokenizer.TokenType() != 3 { // not an identifier
			panic("Unexpected token type! Expected identifier for variable name")
		}
		ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
		ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
		ce.tokenizer.Advance()
		if ce.tokenizer.currentToken == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
			ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
			ce.tokenizer.Advance()
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
	if ce.tokenizer.currentToken != "{" {
		panic("Unexpected token! Expected {")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//3. Loop to handle variable declarations (var):
	//     If the current token is var, call compileVarDec.
	//     Otherwise, break the loop.
	for {
		ce.tokenizer.Advance()
		if ce.tokenizer.currentToken == "var" {
			ce.CompileVarDec()
		} else {
			break //no more vars
		}
	}
	//4. Call compileStatements to handle the statements within the subroutine body.
	ce.tokenizer.Advance()
	ce.CompileStatements()
	//5. Write the closing brace }.
	if ce.tokenizer.currentToken != "}" {
		panic("Unexpected token! Expected }")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//6. Write the closing tag </subroutineBody>.
	ce.WriteCloseTag(ce.hierarchWriter, "subroutineBody")
}

func (ce *CompilationEngine) CompileVarDec() {
	//Purpose: Compiles a var declaration.
	//Steps:
	//1. Write the opening tag <varDec>.
	ce.WriteOpenTag(ce.hierarchWriter, "varDec")
	//2. Write the current token var.
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//3. Advance the tokenizer and write the type.
	if ce.tokenizer.TokenType() != 1 { // not a keyword
		panic("Unexpected token type! Expected keyword for var type")
	}
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//4. Advance the tokenizer and write the first variable name.
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for variable name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
	//5. Loop to handle additional variables:
	//     If the current token is a comma ,, write the comma and the next variable name.
	//     Otherwise, break the loop.
	for {
		ce.tokenizer.Advance()
		if ce.tokenizer.currentToken == "," {
			ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
			ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
			ce.tokenizer.Advance()
			if ce.tokenizer.TokenType() != 3 { // not an identifier
				panic("Unexpected token type! Expected identifier for variable name")
			}
			ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
			ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
		} else {
			break // no more variables
		}
	}
	//6. Write the semicolon ;.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
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
		if ce.tokenizer.currentToken == "let" {
			ce.tokenizer.Advance()
			ce.CompileLet()
		} else if ce.tokenizer.currentToken == "if" {
			ce.tokenizer.Advance()
			ce.CompileIf()
		} else if ce.tokenizer.currentToken == "while" {
			ce.tokenizer.Advance()
			ce.CompileWhile()
		} else if ce.tokenizer.currentToken == "do" {
			ce.tokenizer.Advance()
			ce.CompileDo()
		} else if ce.tokenizer.currentToken == "return" {
			ce.tokenizer.Advance()
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
	ce.WriteXML(ce.hierarchWriter, "keyword", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "keyword", ce.tokenizer.currentToken)
	//3. Advance the tokenizer and write the subroutine name.
	ce.tokenizer.Advance()
	if ce.tokenizer.TokenType() != 3 { // not an identifier
		panic("Unexpected token type! Expected identifier for subroutine name")
	}
	ce.WriteXML(ce.hierarchWriter, "identifier", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "identifier", ce.tokenizer.currentToken)
	//4. Advance the tokenizer and write the opening parenthesis (.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken != "(" {
		panic("Unexpected token! Expected (")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//5. Advance the tokenizer and call compileExpressionList.
	ce.tokenizer.Advance()
	ce.CompileExpressionList()
	//6. Write the closing parenthesis ).
	ce.tokenizer.Advance() //----------need this?
	if ce.tokenizer.currentToken != ")" {
		panic("Unexpected token! Expected )")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//7. Advance the tokenizer and write the semicolon ;.
	ce.tokenizer.Advance()
	if ce.tokenizer.currentToken != ";" {
		panic("Unexpected token! Expected ;")
	}
	ce.WriteXML(ce.hierarchWriter, "symbol", ce.tokenizer.currentToken)
	ce.WriteXML(ce.plainWriter, "symbol", ce.tokenizer.currentToken)
	//8. Write the closing tag </doStatement>.
	ce.WriteCloseTag(ce.hierarchWriter, "doStatement")
}

func (ce *CompilationEngine) CompileLet() {
	/*
		Purpose: Compiles a let statement.
		Steps:
		1. Write the opening tag <letStatement>.
		2. Write the current token let.
		3. Advance the tokenizer and write the variable name.
		4. Advance the tokenizer to check for array indexing:
		      If the current token is an opening bracket [, write the bracket and call compileExpression.
		      Write the closing bracket ] and advance the tokenizer.
		5. Write the equals sign =.
		6. Advance the tokenizer and call compileExpression.
		7. Write the semicolon ;.
		8. Write the closing tag </letStatement>.
	*/
}

func (ce *CompilationEngine) CompileWhile() {
	/*
		Purpose: Compiles a while statement.
		Steps:
		1. Write the opening tag <whileStatement>.
		2. Write the current token while.
		3. Advance the tokenizer and write the opening parenthesis (.
		4. Advance the tokenizer and call compileExpression.
		5. Write the closing parenthesis ).
		6. Advance the tokenizer and write the opening brace {.
		7. Advance the tokenizer and call compileStatements.
		8. Write the closing brace }.
		9. Write the closing tag </whileStatement>.
	*/
}

func (ce *CompilationEngine) CompileReturn() {
	/*
		Purpose: Compiles a return statement.
		Steps:
		1. Write the opening tag <returnStatement>.
		2. Write the current token return.
		3. Advance the tokenizer to check for an expression:
		     If the current token is not a semicolon ;, call compileExpression.
		4. Write the semicolon ;.
		5. Write the closing tag </returnStatement>.
	*/
}

func (ce *CompilationEngine) CompileIf() {
	/*
		Purpose: Compiles an if statement, possibly with a trailing else clause.
		Steps:
		1. Write the opening tag <ifStatement>.
		2. Write the current token if.
		3. Advance the tokenizer and write the opening parenthesis (.
		4. Advance the tokenizer and call compileExpression.
		5. Write the closing parenthesis ).
		6. Advance the tokenizer and write the opening brace {.
		7. Advance the tokenizer and call compileStatements.
		8. Write the closing brace }.
		9. Advance the tokenizer to check for an else clause:
		      If the current token is else, write the keyword else, the opening brace {, call compileStatements, and write the closing brace }.
		10. Write the closing tag </ifStatement>
	*/
}

func (ce *CompilationEngine) CompileExpression() {
	/*
		Purpose: Compiles an expression.
		Steps:
		1. Write the opening tag <expression>.
		2. Call compileTerm.
		3. Loop to handle additional terms connected by operators:
		      If the current token is an operator, write the operator and call compileTerm for the next term.
		      Otherwise, break the loop.
		4. Write the closing tag </expression>
	*/
}

func (ce *CompilationEngine) CompileTerm() {
	/*
		Purpose: Compiles a term.
		Steps:
		1. Write the opening tag <term>.
		2. Depending on the current token, handle different types of terms:
		      If the token is an integer constant, write the integer constant.
		      If the token is a string constant, write the string constant.
		      If the token is a keyword constant, write the keyword.
		      If the token is an identifier, handle variable names, array entries, or subroutine calls.
		      If the token is an opening parenthesis (, call compileExpression and write the closing parenthesis ).
		      If the token is a unary operator, write the operator and call compileTerm.
		3. Write the closing tag </term>.
	*/
}

func (ce *CompilationEngine) CompileExpressionList() {
	/*
		Purpose: The compileExpressionList function is responsible for compiling a (possibly empty) comma-separated list of expressions. This list is typically found within the argument list of a subroutine call.
		Steps:
		1. Write the opening tag <expressionList>.
		2. Check if the current token indicates the start of an expression. This can be identified by looking for tokens that can start an expression such as integer constants, string constants, keyword constants, variable names, subroutine calls, expressions enclosed in parentheses, and unary operators.
		3. If there is at least one expression:
		     Call compileExpression to compile the first expression.
		     Loop to handle additional expressions separated by commas:
		          If the current token is a comma ,, write the comma symbol.
		          Advance the tokenizer.
		          Call compileExpression to compile the next expression.
		4. Write the closing tag </expressionList>.
	*/
}

// helper functions
// - might need to change how it prints
func (ce *CompilationEngine) WriteXML(writer *bufio.Writer, tag string, content string) {
	writer.WriteString(fmt.Sprintf("<%s> %s </%s>\n", tag, content, tag))
}

func (ce *CompilationEngine) WriteOpenTag(writer *bufio.Writer, tag string) {
	writer.WriteString(fmt.Sprintf("<%s>\n", tag))
}

func (ce *CompilationEngine) WriteCloseTag(writer *bufio.Writer, tag string) {
	writer.WriteString(fmt.Sprintf("</%s>\n", tag))
}
