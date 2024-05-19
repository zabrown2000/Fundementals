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
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store sum of D and value in that spot in that space in the stack -> move SP back down
		// to next empty spot
		asm = "//add\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D+M\n@SP\nM=M+1\n"
	case "sub":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store difference of value in that spot and D in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//sub\n@SP\nAM=M-1\nD=M\nA=A-1\nM=M-D\n@SP\nM=M+1\n"
	case "neg":
		// move stack pointer to top of stack -> negate the value in that spot and store
		// it back in that same spot -> move SP back down to next empty spot
		asm = "//neg\n@SP\nA=M-1\nM=-M\n@SP\nM=M+1\n"
	case "and":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store the and-ing of D and value in that spot in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//and\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D&M\n@SP\nM=M+1\n"
	case "or":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store the or-ing of D and value in that spot in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//or\n@SP\nAM=M-1\nD=M\nA=A-1\nM=D|M\n@SP\nM=M+1\n"
	case "not":
		// move stack pointer to top of stack -> not the value in that spot and store
		// it back in that same spot -> move SP back down to next empty spot
		asm = "//not\n@SP\nA=M-1\nM=!M\n@SP\nM=M+1\n"
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
func (cw *CodeWriter) WritePushPop(command string, segment string, index int) {
	var asm string
	index_str := strconv.Itoa(index)
	if command == "push" {
		cw.writer.Write([]byte(" //push " + segment + " " + index_str))
		switch segment {
		case "constant":
			// put value to go in stack in A and then D -> set A to next open spot in stack and place
			// value there -> move SP down 1 to next open spot
			asm = "@" + index_str + "\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "local":
			// get address of local segment and store segment in D -> get address in local segment
			// based on index offset -> store value in that spot in D -> place value from D in next
			// open spot in stack -> move SP down 1 to next open spot
			asm = "@LCL\nD=M\n@" + index_str + "\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "argument":
			// get address of argument segment and store segment in D -> get address in argument segment
			// based on index offset -> store value in that spot in D -> place value from D in next
			// open spot in stack -> move SP down 1 to next open spot
			asm = "@ARG\nD=M\n@" + index_str + "\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "this":
			// get address of this segment and store segment in D -> get address in this segment
			// based on index offset -> store value in that spot in D -> place value from D in next
			// open spot in stack -> move SP down 1 to next open spot
			asm = "@THIS\nD=M\n@" + index_str + "\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "that":
			// get address of that segment and store segment in D -> get address in that segment
			// based on index offset -> store value in that spot in D -> place value from D in next
			// open spot in stack -> move SP down 1 to next open spot
			asm = "@THAT\nD=M\n@" + index_str + "\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "pointer":
			// based on index decide if taking this or that segment address -> store the segment in D
			// -> place segment in next open spot in stack -> move SP down 1 to next open spot
			if index == 0 {
				asm = "@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
			} else if index == 1 {
				asm = "@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
			}
		case "temp":
			// get address of beginning of temp registers and move down to correct one based on index
			// -> store value of that register in D -> get address in that segment
			// based on index offset -> store value in that spot in D -> place value from D in next
			// open spot in stack -> move SP down 1 to next open spot
			asm = "@R5\nD=A\n@" + index_str + "\nA=D+A\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		case "static":
			// determine class name and # static field we are up to -> A becomes the index+1-th register
			// address in the static segment -> store value in that register in D -> place value from D
			// in next open spot in stack -> move SP down 1 to next open spot
			asm = "@" + cw.file_name + "." + index_str + "\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n"
		}
	} else if command == "pop" {
		//insert pop stuff here
		cw.writer.Write([]byte(" //pop " + segment + " " + index_str))
	}
	cw.writer.Write([]byte(asm))
	cw.writer.Flush()
}
