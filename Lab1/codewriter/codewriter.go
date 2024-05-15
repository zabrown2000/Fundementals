package codewriter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// TO DO: refactor all these functions to output asm code
// add in and, not, or
// 5 functions to add from textbook - ctor, setFileName, writeArithmatic, writePushPop, close
// Zahava: add, sub, neg, and, or, not, push
// Tali: eq, gt, lt, pop

type CodeWriter struct {
	writer    *bufio.Writer
	file_name string
}

func New(asm_path string) *CodeWriter {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("File opening error", err)
		return nil
	}
	return &CodeWriter{bufio.NewWriter(write_file), ""}
}

func (cw *CodeWriter) SetFileName(name string) {
	cw.file_name = name
}

// TO DO: once all written together, run in emulator to see each step to make sure understand
func (cw *CodeWriter) WriteArithmetic(cmd string) {
	var asm string
	switch cmd {
	case "add":
		asm = "//add\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D+M\n"
	case "sub":
		asm = "//sub\n@SP\nAM=M-1\nD=M\nA=A-1\nM=M-D\n"
	case "neg":
		asm = "//neg\n@SP\nA=M-1\nM=-M"
	case "and":
		asm = "//and\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D&M\n"
	case "or":
		asm = "//or\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D|M\n"
	case "not":
		asm = "//not\n@SP\nA=M-1\nM=!M\n"
	case "eq":
	//tali
	case "gt":
	//tali
	case "lt":
		//tali
	}

	cw.writer.Write([]byte(asm))
	cw.writer.Flush()
}

// maybe refactor and make subfunctions
func (cw *CodeWriter) writePushPop(command string, segment string, index int) {
	var asm string
	index_str := strconv.Itoa(index)
	if command == "C_PUSH" {
		switch segment {
		case "constant":
			asm = "@" + index_str + "\nD=A\n"
		case "local":
			asm = "@LCL\nD=M\n@" + index_str + "\nA=D+A\nD=M\n"
		case "argument":
			asm = "@ARG\nD=M\n@" + index_str + "\nA=D+A\nD=M\n"
		case "this":
			asm = "@THIS\nD=M\n@" + index_str + "\nA=D+A\nD=M\n"
		case "that":
			asm = "@THAT\nD=M\n@" + index_str + "\nA=D+A\nD=M\n"
		case "pointer":
			if index == 0 {
				asm = "@THIS\nD=M\n"
			} else {
				asm = "@THAT\nD=M\n"
			}
		case "static":
			asm = "@" + cw.file_name + "." + index_str + "\nD=M\n"
		case "temp":
			asm = "@R5\nD=A\n@" + index_str + "\nA=D+A\nD=M\n"
		}
	} else if command == "C_POP" {
		//insert pop stuff here
	}
	cw.writer.Write([]byte(asm))
	cw.writer.Flush()
}

/*
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


func handlePop(str string, num int) {
	writer, err := createWriter()
	if err != nil {
		fmt.Println("Error creating writer:", err)
		return
	}
	writer.Write([]byte("command: pop segment: " + str + " index: " + strconv.Itoa(num) + "\n"))
	writer.Flush()
}*/
