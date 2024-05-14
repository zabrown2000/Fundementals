package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	// codewriter/codewriter
	// parser/parser
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

	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	defer write_file.Close()

	//maybe here create codewriter and set filename to be asm_file_name

	writer := bufio.NewWriter(write_file)

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".vm" {
			CurVM = file.Name()
			curFile, err := os.Open(dir_path + CurVM)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer curFile.Close()
			reader := bufio.NewReader(curFile)
			counter = 1

			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				words := strings.Fields(line)
				if len(words) != 0 {
					switch words[0] {
					case "add":
						handleAdd()
					case "sub":
						handleSub()
					case "neg":
						handleNeg()
					case "eq":
						handleEq()
					case "gt":
						handleGt()
					case "lt":
						handleLt() // add in and,or,not
					case "push":
						s, err := strconv.Atoi(words[2])
						if err != nil {
							fmt.Println("Can't convert this to an int!")
						} else {
							handlePush(words[1], s)
						}
					case "pop":
						s, err := strconv.Atoi(words[2])
						if err != nil {
							fmt.Println("Can't convert this to an int!")
						} else {
							handlePop(words[1], s)
						}
					default:
						writer.Write([]byte("unknown\n"))
						writer.Flush()
					}
				}
			}
			fmt.Println("End of input file: " + file.Name())
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
