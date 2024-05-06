package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
)

var CurVM string

// path: C:\Users\zbrow\OneDrive\Documents\Machon_Tal\Fundementals\Lab0-Eng\

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	fmt.Scanln(&dir_path)
	var asm_path = dir_path + "Lab0-Eng.asm"

	// create output file Tar0.asm in that directory and open from writing text
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse files in dir with extention .vm
	for _, file := range files {
		if path.Ext(file.Name()) == ".vm" {
			fmt.Println(file.Name())
			CurVM = file.Name()
			curFile, err := os.Open(dir_path + CurVM)
			if err != nil {
				fmt.Println(err)
				return
			}

			reader := bufio.NewReader(curFile)
			var counter int = 1

			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				// take out 1st word, swicth it, if add send it, if logical inc ctr and then send with ctr, if mem get rest of words and send
				// or have helper fn inc ctr and have global ctr and then reset it for each file
				//fmt.Println(line)
				str := strconv.Itoa(counter)
				str = str + ":"
				writer.Write([]byte(str + line))
				writer.Flush()
				counter++
			}
		}
	}

	// each file hs own loop
	// in loop have ctr for num of logical commands
	// have global var for name of file without .vm
	// read each line to decide helper funciton - switch
	// at end of ech file close file and print to screen
	// after outer loop make print
	// write the helper functions
	// handleAdd, handleSub, handleNeg
	// handleEq, HandleGt, handleLt
	// handlePush, handlePop
}
