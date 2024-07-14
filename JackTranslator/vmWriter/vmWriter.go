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
