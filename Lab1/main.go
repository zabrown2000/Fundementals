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
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}
	dir_name = filepath.Base(dir_path)
	asm_file_name = dir_name + ".asm"
	asm_path = dir_path + asm_file_name

	// create codewriter obj and send file to open to write
	cw := codewriter.New(asm_path)

	//writer := bufio.NewWriter(write_file) - get from codewriter

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// create loop for dir

	// if arithmetic get cmd type from arg1, and call writearithmetic
	// if push or pop call the push/pop func and send arg1 and arg2
	for _, file := range files {
		if path.Ext(file.Name()) == ".vm" {
			// call setfilename and send file name without vm - basically dir_name
			// fileNameWithoutExt := strings.TrimSuffix(filePath, ".vm")
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
				switch cmdType {
				case parser.C_ARITHMETIC:
					//get arg1 and send in below func
					fmt.Println("C_ARITHMETIC")
					arg1 := psr.Arg1()
					cw.WriteArithmetic(arg1)
				case parser.C_PUSH:
					//get arg1 and arg2 and send
					fmt.Println("C_PUSH")
					arg1 := psr.Arg1()
					arg2 := psr.Arg2()
					cw.WritePushPop("push", arg1, arg2)
				case parser.C_POP:
					//get arg1 and arg2 and send
					fmt.Println("C_POP")
					arg1 := psr.Arg1()
					arg2 := psr.Arg2()
					cw.WritePushPop("pop", arg1, arg2)
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
