// Zahava Brown - 557029367
// Tali Cohen - 329651871
package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

// path: C:\Users\zbrow\OneDrive\Documents\Machon_Tal\Fundementals\Lab0-Eng\
// path: C:\Users\Merekat\Documents\School\23-24\Fundamentals\Lab0-Eng\

var CurVM string

var counter int
var asm_path string

func createWriter() (*bufio.Writer, error) {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil, err
	}
	writer := bufio.NewWriter(write_file)
	return writer, nil
}

func handleAdd() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: add\n"))
	writer.Flush()
}

func handleSub() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: sub\n"))
	writer.Flush()
}

func handleNeg() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: neg\n"))
	writer.Flush()
}

func handleEq() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: eq\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handleGt() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: gt\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handleLt() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: lt\n"))
	writer.Flush()
	str := strconv.Itoa(counter)
	str = str + "\n"
	counter++
	writer.Write([]byte("counter: " + str))
	writer.Flush()
}

func handlePush(str string, num int) {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: push segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	writer.Flush()
}

func handlePop(str string, num int) {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: pop segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	writer.Flush()
}

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	fmt.Scanln(&dir_path)
	var asm_file_name = "Lab0-Eng.asm"
	asm_path = dir_path + asm_file_name

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
			fmt.Println("End of input file: " + file.Name())
		}
	}
	fmt.Println("Output file is ready: " + asm_file_name)
}
