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
	_, err = writer.Write([]byte("command: add\n"))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func handleSub() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	_, err = writer.Write([]byte("command: sub\n"))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func handleNeg() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	_, err = writer.Write([]byte("command: neg\n"))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func handleEq() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	str := strconv.Itoa(counter)
	str = str + "\n"
	_, err = writer.Write([]byte("command: eq\ncounter: " + str))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	counter++
}

func handleGt() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	str := strconv.Itoa(counter)
	str = str + "\n"
	_, err = writer.Write([]byte("command: gt\ncounter: " + str))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	counter++
}

func handleLt() {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	str := strconv.Itoa(counter)
	str = str + "\n"
	_, err = writer.Write([]byte("command: lt\ncounter: " + str))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	counter++
}

func handlePush(str string, num int) {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	_, err = writer.Write([]byte("command: push segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func handlePop(str string, num int) {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	_, err = writer.Write([]byte("command: pop segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}
	var asm_file_name = "Lab0-Eng.asm"
	asm_path = dir_path + asm_file_name

	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return
	}
	defer func(write_file *os.File) {
		err := write_file.Close()
		if err != nil {

		}
	}(write_file)
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
			defer func(curFile *os.File) {
				err := curFile.Close()
				if err != nil {

				}
			}(curFile)
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
						_, err := writer.Write([]byte("unknown\n"))
						if err != nil {
							return
						}
						err = writer.Flush()
						if err != nil {
							return
						}
					}
				}
			}
			fmt.Println("End of input file: " + file.Name())
		}
	}
	fmt.Println("Output file is ready: " + asm_file_name)
}
