package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

var CurVM string

// var writer *bufio.Writer
var counter int
var asm_path string

func handleAdd() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: add\n"))
	writer.Flush()
}

func handleSub() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: sub\n"))
	writer.Flush()
}

func handleNeg() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: neg\n"))
	writer.Flush()
}

func handleEq() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: eq\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handleGt() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: gt\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handleLt() {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: lt\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handlePush(str string, num int) {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: push segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	writer.Flush()
}

func handlePop(str string, num int) {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	writer := bufio.NewWriter(write_file)
	writer.Write([]byte("command: pop segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	writer.Flush()
}

// path: C:\Users\zbrow\OneDrive\Documents\Machon_Tal\Fundementals\Lab0-Eng\
// path: C:\Users\Merekat\Documents\School\23-24\Fundamentals\Lab0-Eng\

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	fmt.Scanln(&dir_path)
	asm_path = dir_path + "Lab0-Eng.asm"

	// create output file Tar0.asm in that directory and open from writing text
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	defer write_file.Close()
	writer := bufio.NewWriter(write_file)

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse files in dir with extension .vm
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
				// take out 1st word, switch it, if add send it, if logical inc ctr and then send with ctr, if mem get rest of words and send
				// or have helper fn inc ctr and have global ctr and then reset it for each file
				//fmt.Println(line)
				//str := strconv.Itoa(counter)
				//str = str + ":"
				//writer.Write([]byte(str + line))
				//writer.Flush()
				//counter++
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
						handleLt()
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
			fmt.Println("End of input file: " + file.Name() + "\n")
		}
	}
	fmt.Println("Output file is ready: " + write_file.Name() + "\n")
	// each file hs own loop
	// in loop have ctr for num of logical commands
	// have global var for name of file without .vm
	// read each line to decide helper function - switch - done
	// at end of ech file close file and print to screen
	// after outer loop make print
	// write the helper functions - all done
	// handleAdd, handleSub, handleNeg
	// handleEq, HandleGt, handleLt
	// handlePush, handlePop
}
