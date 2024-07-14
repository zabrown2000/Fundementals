package vmWriter

import (
	"fmt"
	"os"
)

/*
constructor(file stream)
writePush(segment)
writePop(segment)
writeArithmetic(command)
writeLabel(label)
WriteGoto(label)
writeIf(label)
writeCall(name, nargs)
writeFunction(name, nlocals)
writeReturn
*/

func push(file *os.File, segment string, index int) {
	_, err := fmt.Fprintf(file, "push %s %d\n", segment, index)
	if err != nil {
		return
	}
}

func pop(file *os.File, segment string, index int) {
	_, err := fmt.Fprintf(file, "pop %s %d\n", segment, index)
	if err != nil {
		return
	}
}

func label(file *os.File, labelName string) {
	_, err := fmt.Fprintf(file, "label %s\n", labelName)
	if err != nil {
		return
	}
}

func goTo(file *os.File, labelName string) {
	_, err := fmt.Fprintf(file, "goto %s\n", labelName)
	if err != nil {
		return
	}
}

func ifGoto(file *os.File, labelName string) {
	_, err := fmt.Fprintf(file, "if-goto %s\n", labelName)
	if err != nil {
		return
	}
}

func function(file *os.File, functionName string, nVars int) {
	_, err := fmt.Fprintf(file, "function %s %d\n", functionName, nVars)
	if err != nil {
		return
	}
}

func call(file *os.File, functionName string, nArgs int) {
	_, err := fmt.Fprintf(file, "call %s %d\n", functionName, nArgs)
	if err != nil {
		return
	}
}

func returnFunc(file *os.File) {
	_, err := fmt.Fprintf(file, "return\n")
	if err != nil {
		return
	}
}
