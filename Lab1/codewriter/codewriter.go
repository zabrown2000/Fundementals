package codewriter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// TO DO: refactor all these functions to output asm code
// add in and, not, or
// 5 functions to add from textbook
// Zahava: add, sub, neg, and, or, not, push
// Tali: eq, gt, lt, pop

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
	str := strconv.Itoa(counter)
	str = str + "\n"
	writer.Write([]byte("command: eq\ncounter: " + str))
	writer.Flush()
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
	writer.Write([]byte("command: gt\ncounter: " + str))
	writer.Flush()
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
	writer.Write([]byte("command: lt\ncounter: " + str))
	writer.Flush()
	counter++
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
