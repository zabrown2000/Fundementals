package codewriter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type CodeWriter struct {
	writer              *bufio.Writer
	file_name           string
	logic_counter       int
	vm_function_counter int
	// convention: filemame.funcname.counter
	// vm file already has filename.funcname
	// when need ctr after label?
	// to allow mult calls to same func in code, when write call in asm
	// first have to store return label and return takes that address to go back to
	// label to go and come back need to be the same
}

// New creates a new CodeWriter instance for a given file path
func New(asm_path string) *CodeWriter {
	write_file, err := os.OpenFile(asm_path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	logic_counter := 0
	vm_function_counter := 0
	if err != nil {
		fmt.Println("File opening error", err)
		return nil
	}
	return &CodeWriter{bufio.NewWriter(write_file), "", logic_counter, vm_function_counter}
}

// SetFileName sets the file name for the current vm file for dealing with static segment
func (cw *CodeWriter) SetFileName(name string) {
	cw.file_name = name
}

// WriteInit writes the bootstrap code at the top of the asm file
func (cw *CodeWriter) WriteInit() {
	//func (cw *CodeWriter) WriteInit(callSysInit bool) {
	//	cw.WriteLine("@256\nD=A\n@SP\nM=D\n")
	//	// if sys.vm provided call Sys.init 0
	//	if callSysInit {
	//		cw.WriteComment("call Sys.init 0")
	//		cw.WriteCall("Sys.init", 0)
	//		cw.WriteLine("0;JMP\n")
	//	}
	//}

	// Need to save all segments when call sysinit
	// Set stack pointer to 256
	var asm = "//set stack pointer\n@256\nD=A\n@0\nM=D\n"
	// call writeCall, it will save segments for me
	//@Sys.init$ret.0\nD=A\n@SP\nA=M\nM=D\n@SP\nM=M+1\n - push return address
	//@LCL\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n
	//@ARG\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n
	//@THIS\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n
	//@THAT\nD=M\n@SP\nA=M\nM=D\n@SP\nM=M+1\n
	//@SP\nD=M\n@5\nD=D-A\n@ARG\nM=D\n - reposition ARG
	//@SP\nD=M\n@LCL\nM=D\n - reposition LCL
	asm += "@Sys.init\n0;JMP\n(Sys.init$ret.0)\n" //go to sysinit

	_, err := cw.writer.Write([]byte(asm))
	if err != nil {
		return
	}
	err = cw.writer.Flush()
	if err != nil {
		return
	}
}

// WriteArithmetic writes the assembly code for the given VM arithmetic command
func (cw *CodeWriter) WriteArithmetic(cmd string) {
	var asm string
	switch cmd {
	case "add":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store sum of D and value in that spot in that space in the stack -> move SP back down
		// to next empty spot
		asm = "//add\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D+M\n@SP\nM=M+1\n"
	case "sub":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store difference of value in that spot and D in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//sub\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=M-D\n@SP\nM=M+1\n"
	case "neg":
		// move stack pointer to top of stack -> negate the value in that spot and store
		// it back in that same spot -> move SP back down to next empty spot
		asm = "//neg\n@SP\nAM=M-1\nM=-M\n@SP\nM=M+1\n"
	case "and":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store the and-ing of D and value in that spot in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//and\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D&M\n@SP\nM=M+1\n"
	case "or":
		// move stack pointer to top of stack -> store value in D -> move up to the next element
		// -> store the or-ing of D and value in that spot in that space in the stack -> move SP back
		// down to next empty spot
		asm = "//or\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nM=D|M\n@SP\nM=M+1\n"
	case "not":
		// move stack pointer to top of stack -> note the value in that spot and store
		// it back in that same spot -> move SP back down to next empty spot
		asm = "//not\n@SP\nAM=M-1\nM=!M\n@SP\nM=M+1\n"
	case "eq":
		// move SP to top of stack and store value in D
		//decrease the SP and subtract the value in D from the contents in M and store in D
		//if D = 0 jump to label IS-EQ, load the SP and set contents of stack at SP to -1 (true)
		// else load SP and set content of stack to 0 (false) -> move SP back down to next empty spot
		asm = "//equal\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D\n@IS_EQ_LABEL" + strconv.Itoa(cw.logic_counter) + "\nD;JEQ\n//Not_Equal\n@SP\nA=M\nM=0\n@END_EQ_LABEL" + strconv.Itoa(cw.logic_counter) + "\n0;JMP\n(IS_EQ_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nA=M\nM=-1\n(END_EQ_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nM=M+1\n"
		cw.logic_counter++
	case "gt":
		// move SP to top of stack and store value in D
		//decrease the SP and subtract the value in D from the contents in M and store in D
		//if D > 0 jump to label IS-GT, load the SP and set contents of stack at SP to -1 (true)
		// else load SP and set content of stack to 0 (false) -> move SP back down to next empty spot
		asm = "//gt\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D\n@IS_GT_LABEL" + strconv.Itoa(cw.logic_counter) + "\nD;JGT\n//Not_GT\n@SP\nA=M\nM=0\n@END_GT_LABEL" + strconv.Itoa(cw.logic_counter) + "\n0;JMP\n(IS_GT_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nA=M\nM=-1\n(END_GT_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nM=M+1\n"
		cw.logic_counter++
	case "lt":
		// move SP to top of stack and store value in D
		//decrease the SP and subtract the value in D from the contents in M and store in D
		//if D < 0 jump to label IS-LT, load the SP and set contents of stack at SP to -1 (true)
		// else load SP and set content of stack to 0 (false) -> move SP back down to next empty spot
		asm = "//lt\n@SP\nAM=M-1\nD=M\n@SP\nAM=M-1\nD=M-D\n@IS_LT_LABEL" + strconv.Itoa(cw.logic_counter) + "\nD;JLT\n//Not_LT\n@SP\nA=M\nM=0\n@END_LT_LABEL" + strconv.Itoa(cw.logic_counter) + "\n0;JMP\n(IS_LT_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nA=M\nM=-1\n(END_LT_LABEL" + strconv.Itoa(cw.logic_counter) + ")\n@SP\nM=M+1\n"
		cw.logic_counter++
	}

	_, err := cw.writer.Write([]byte(asm))
	if err != nil {
		return
	}
	err = cw.writer.Flush()
	if err != nil {
		return
	}
}

// WritePushPop writes the assembly code for the given VM push/pop command
func (cw *CodeWriter) WritePushPop(command string, segment string, index int) {
	var asm string
	index_str := strconv.Itoa(index)
	if command == "push" {
		_, err := cw.writer.Write([]byte("//push " + segment + " " + index_str + "\n"))
		if err != nil {
			return
		}
		err = cw.writer.Flush()
		if err != nil {
			return
		}
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
		_, err := cw.writer.Write([]byte("//pop " + segment + " " + index_str + "\n"))
		if err != nil {
			return
		}
		err = cw.writer.Flush()
		if err != nil {
			return
		}
		switch segment {
		case "constant":
			// Load SP and decrease to access top value, and store it in D
			asm = "@SP\nM=M-1\nD=M\n"
		case "local":
			// get address of local segment and store segment in D -> calculate internal address in local segment
			// by adding index offset -> store address in General Register 13
			//Load SP and decrement to access top value and store in D
			//Load address of R13 and access it contents (local address with offset)
			//store contents in D (top value from stack) in to that address
			asm = "@LCL\nD=M\n@" + index_str + "\nD=D+A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		case "argument":
			// get address of argument segment and store segment in D -> calculate internal address in argument segment
			// by adding index offset -> store address in General Register 13
			//Load SP and decrement to access top value and store in D
			//Load address of R13 and access it contents (argument address with offset)
			//store contents in D (top value from stack) in to that address
			asm = "@ARG\nD=M\n@" + index_str + "\nD=D+A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		case "this":
			// get address of 'this' segment and store segment in D -> calculate internal address in 'this' segment
			// by adding index offset -> store address in General Register 13
			//Load SP and decrement to access top value and store in D
			//Load address of R13 and access it contents ('this' address with offset)
			//store contents in D (top value from stack) in to that address
			asm = "@THIS\nD=M\n@" + index_str + "\nD=D+A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		case "that":
			// get address of 'that' segment and store segment in D -> calculate internal address in 'that' segment
			// by adding index offset -> store address in General Register 13
			//Load SP and decrement to access top value and store in D
			//Load address of R13 and access it contents ('that' address with offset)
			//store contents in D (top value from stack) in to that address
			asm = "@THAT\nD=M\n@" + index_str + "\nD=D+A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		case "pointer":
			// based on index decide if storing this or that segment address -> reduce the SP to top value in stack
			//load that value to D and load THIS/THAT address to and store contents of D there
			if index == 0 {
				asm = "@SP\nAM=M-1\nD=M\n@THIS\nM=D\n"
			} else if index == 1 {
				asm = "@SP\nAM=M-1\nD=M\n@THAT\nM=D\n"
			}
		case "temp":
			// get address of beginning of temp registers and move down to correct one based on index
			// calculate the new address in D, then store it in R13,
			// Load SP, decrease it to pop top value into D,
			// reload temp register address from R13 and store D at the corresponding address
			asm = "@R5\nD=A\n@" + index_str + "\nD=D+A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		case "static":
			// determine class name and # static field we are up to -> A becomes the index+1-th register
			// address in the static segment -> store that address in R13
			// Load SP, decrease it to pop top value into D,
			// reload static address from R13 and store value in D there
			asm = "@" + cw.file_name + "." + index_str + "\nD=A\n@R13\nM=D\n@SP\nAM=M-1\nD=M\n@R13\nA=M\nM=D\n"
		}
	}
	_, err := cw.writer.Write([]byte(asm))
	if err != nil {
		return
	}
	err = cw.writer.Flush()
	if err != nil {
		return
	}
}

// WriteLabel writes the assembly code for the labelprovided
func (cw *CodeWriter) WriteLabel(label string) {
	var asm = "//label\n(" + label + ")\n"

	_, err := cw.writer.Write([]byte(asm))
	if err != nil {
		return
	}
	err = cw.writer.Flush()
	if err != nil {
		return
	}
}

// WriteGoto writes the assembly code for an unconditional jump
func (cw *CodeWriter) WriteGoto(label string) {
	var asm = "//goto\n@" + label + "\n0;JMP\n"

	_, err := cw.writer.Write([]byte(asm))
	if err != nil {
		return
	}
	err = cw.writer.Flush()
	if err != nil {
		return
	}
}

func (cw *CodeWriter) WriteIfGoto(label string) {
	// TO DO: Tali - add asm code and write it
	//               - update code in main to handle this
}

func (cw *CodeWriter) WriteFunction(function_name string, nVars int) {
	// TO DO: Tali - add asm code and write it
	//               - update code in main to handle this
}

func (cw *CodeWriter) WriteReturn(function_name string, nVars int) {
	// TO DO: Tali - add asm code and write it
	//               - update code in main to handle this
}

func (cw *CodeWriter) WriteCall(function_name string, nArgs int) {
	// TO DO: Zahava - add asm code and write it
	//               - update code in main to handle this

	// create label for return address - filename.funcname$ret.ctr (put $ret for return addresses)
	// push the return address with the label
	// save LCL, ARG, THIS, THAT segments
	// reposition ARG
	// reposition LCL
	// place label of function going to and jump to it

	//func (cw *CodeWriter) WriteCall(funcName string, nArgs int) {
	//	count := cw.nextCount()
	//	cw.WriteLine("@SP\nD=M\n@R13\nM=D\n" +
	//		"@ret." + count + "\nD=A\n@SP\nA=M\nM=D\n" +
	//		tplIncrementSP() +
	//		tplPointer("LCL") +
	//		tplIncrementSP() +
	//		tplPointer("ARG") +
	//		tplIncrementSP() +
	//		tplPointer("THIS") +
	//		tplIncrementSP() +
	//		tplPointer("THAT") +
	//		tplIncrementSP() +
	//		"@R13\nD=M\n@" + strconv.Itoa(nArgs) + "\nD=D-A\n@ARG\nM=D\n" +
	//		"@SP\nD=M\n@LCL\nM=D\n@" + funcName + "\n" +
	//		"0;JMP\n(ret." + count + ")\n")
	//}
	//
	//func tplIncrementSP() string {
	//	return "@SP\nM=M+1\n"
	//}
	//
	//func tplPointer(i string) string {
	//	return "@" + i + "\nD=M\n@SP\nA=M\nM=D\n"
	//}
}
