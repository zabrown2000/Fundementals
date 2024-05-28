package main

import (
	"Fundementals/VMTranslator/codewriter"
	"Fundementals/VMTranslator/parser"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var CurVM string
var asm_path string
var dir_name string
var asm_file_name string

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}
	dir_name = filepath.Base(dir_path)
	asm_file_name = dir_name + ".asm"
	asm_path = dir_path + asm_file_name

	// create codewriter obj and send file to open to write
	cw := codewriter.New(asm_path)
	// Calling bootstrap code writer
	cw.WriteInit()

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".vm" {
			CurVM = file.Name()
			cw.SetFileName(strings.TrimSuffix(CurVM, ".vm"))
			// each vm file, create parser obj and send file to open to read
			psr, err := parser.New(dir_path + CurVM)
			if err != nil {
				fmt.Println(err)
				return
			}
			for {
				// for each file call parser to read the line and return type of command and args
				cmdType := psr.CommandType()
				// TO DO: Tali & Zahava - add cases for commands added
				switch cmdType {
				case parser.C_ARITHMETIC:
					arg1 := psr.Arg1()
					cw.WriteArithmetic(arg1)
				case parser.C_PUSH:
					arg1 := psr.Arg1()
					arg2 := psr.Arg2()
					cw.WritePushPop("push", arg1, arg2)
				case parser.C_POP:
					arg1 := psr.Arg1()
					arg2 := psr.Arg2()
					cw.WritePushPop("pop", arg1, arg2)
				case parser.C_LABEL:
					arg1 := psr.Arg1()
					cw.WriteLabel(arg1)
				case parser.C_GOTO:
					arg1 := psr.Arg1()
					cw.WriteGoto(arg1)
				case -1: //not a valid commandtype returned
					if psr.HasMoreLines() { //but still more lines
						psr.Advance()
						continue
					}
				default:
					panic("unhandled default case")
				}
				// check if more lines and then if not break
				if psr.HasMoreLines() {
					psr.Advance()
				} else {
					break
				}
			}
			fmt.Println("End of input file: " + file.Name())
		}
	}
	fmt.Println("Output file is ready: " + asm_file_name)
}
