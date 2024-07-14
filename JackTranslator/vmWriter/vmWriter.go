package vmWriter

import (
	"bufio"
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

type VMWriter struct {
	writer    *bufio.Writer
	file_name string
}

// New creates a new VMWriter instance for a given file path
func New(vm_path string) *VMWriter {
	write_file, err := os.OpenFile(vm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil
	}
	return &VMWriter{bufio.NewWriter(write_file), ""}
}

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

func arithmetic(file *os.File, command string) {
	validCommands := []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}
	for _, validCommand := range validCommands {
		if command == validCommand {
			_, err := fmt.Fprintf(file, "%s\n", command)
			if err != nil {
				return
			}
			return
		}
	}
}

func jackFunction(file *os.File, className, functionName string, nVars int) {
	_, err := fmt.Fprintf(file, "function %s.%s %d\n", className, functionName, nVars)
	if err != nil {
		return
	}
}

func jackMethod(file *os.File, className, methodName string, nVars int) {
	methodDeclaration := fmt.Sprintf("function %s.%s %d\n", className, methodName, nVars)
	setupThisPointer := "push argument 0\npop pointer 0\n"
	_, err := fmt.Fprintf(file, "%s%s", methodDeclaration, setupThisPointer)
	if err != nil {
		return
	}
}

func jackConstructor(file *os.File, className string, nVars, size int) {
	constructorDeclaration := fmt.Sprintf("function %s.new %d\n", className, nVars)
	allocateMemory := fmt.Sprintf("push constant %d\ncall Memory.alloc 1\npop pointer 0\n", size)
	_, err := fmt.Fprintf(file, "%s%s", constructorDeclaration, allocateMemory)
	if err != nil {
		return
	}
}
