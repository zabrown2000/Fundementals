package parser

//parser
/*
import(
	"bufio"
	"os"
	"strings"
	"strconv"

)
*/
//Tali: ctor, hasMoreCommands, advance, commandType, ar1, arg2
//Question - test files have lines that are started with '//' how to ignore them?

//constructor

func hasMoreCommands() bool {
	return false
}

func advance() {
	//moves to next command
}

// Command types
const (
	C_ARITHMETIC = iota // arithmetic (add,sub,neg,eq,lt,gt,and,or,not)
	C_PUSH              // push
	C_POP               // pop
	C_LABEL             // label
	C_GOTO              // goto
	C_IF                // if
	C_FUNCTION          // function
	C_RETURN            // return
	C_CALL              // call
)

func commandType(cmd string) int {
	return C_CALL
}
