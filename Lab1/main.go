package main

import (
	"Fundementals/Lab1/codewriter"
	"Fundementals/Lab1/parser"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//TO DO: clean up main to sync with 2 new modules

var CurVM string
var asm_path string
var dir_name string
var asm_file_name string

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	fmt.Scanln(&dir_path)
	dir_name = filepath.Base(dir_path)
	asm_file_name = dir_name + ".asm"
	asm_path = dir_path + asm_file_name

	//TO DO: zahava - create codewriter object and call setfilename and send current vm file name
	//       zahava - create shell code for entire process using parser and codewriter
	// in loop for each vm file, each vm file call setFileName and send name without .vm attached

	// create codewriter obj and send file to open to write
	cw := codewriter.New(asm_path)

	//writer := bufio.NewWriter(write_file) - get from codewriter

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// create loop for dir

	// if arith get cmd type from arg1, and call writearith
	// if push or pop call the pushpop func and send arg1 and arg2
	for _, file := range files {
		if path.Ext(file.Name()) == ".vm" {
			// call setfilename and send file name without vm - basically dir_name
			// fileNameWithoutExt := strings.TrimSuffix(filePath, ".vm")
			cw.SetFileName(dir_name)
			CurVM = file.Name()
			cw.SetFileName(strings.TrimSuffix(CurVM, ".vm"))
			// each vm file, create parser obj and send file to open to read

			for {
				// for each file call parser to read the line and return type of command and args
				cmdType = parser.commandType(...)
				switch cmdType {
				case parser.C_ARITHMETIC:
					//get arg1 and send in below func
					cw.WriteArithmetic()
				case parser.C_PUSH:
					//get arg1 and arg2 and send
					cw.WritePushPop("push")
				case parser.C_POP:
					//get arg1 and arg2 and send
					cw.WritePushPop("pop")
				}
			}
		}
	}
	fmt.Println("Output file is ready: " + asm_file_name)
}

/*
Testing Order:
1. SimpleAdd: This program pushes two constants onto the stack, and adds them up.
              Tests how your implementation handles the commands “push constant i”, and “add”.

2. StackTest: Pushes some constants onto the stack, and tests how your implementation handles all
              the VM arithmetic-logical commands.
              - vm file contains push constant, eq, lg, gt, add, neg, and, or, not

3. BasicTest: Executes push, pop, and arithmetic commands using the memory segments constant,
              local, argument, this, that, and temp. Tests how your implementation handles these memory
              segments (you've already handled constant).

4. PointerTest: Executes push, pop, and arithmetic commands using the memory segments pointer,
                this, and that. Tests how your implementation handles the pointer segment.

5. StaticTest: Executes push, pop, and arithmetic commands using constants and the memory
               segment static. Tests how your implementation handles the static segment.
*/
