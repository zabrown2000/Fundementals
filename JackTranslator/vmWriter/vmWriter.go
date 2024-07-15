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

func (vm *VMWriter) WritePush(segment string, index int) {
	_, err := fmt.Fprintf(vm.writer, "push %s %d\n", segment, index)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WritePop(segment string, index int) {
	_, err := fmt.Fprintf(vm.writer, "pop %s %d\n", segment, index)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteLabel(labelName string) {
	_, err := fmt.Fprintf(vm.writer, "label %s\n", labelName)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteGoTo(labelName string) {
	_, err := fmt.Fprintf(vm.writer, "goto %s\n", labelName)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteIfGoto(labelName string) {
	_, err := fmt.Fprintf(vm.writer, "if-goto %s\n", labelName)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteFunction(functionName string, nVars int) {
	_, err := fmt.Fprintf(vm.writer, "function %s %d\n", functionName, nVars)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteCall(functionName string, nArgs int) {
	_, err := fmt.Fprintf(vm.writer, "call %s %d\n", functionName, nArgs)
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteReturn() {
	_, err := fmt.Fprintf(vm.writer, "return\n")
	if err != nil {
		return
	}
}

func (vm *VMWriter) WriteArithmetic(command string) {
	validCommands := []string{"add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not"}
	for _, validCommand := range validCommands {
		if command == validCommand {
			_, err := fmt.Fprintf(vm.writer, "%s\n", command)
			if err != nil {
				return
			}
			return
		}
	}
}
