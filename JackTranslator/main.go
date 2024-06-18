package JackTranslator

import (
	"fmt"
	"os"
	"path"
)

// add catch for if parser throws failure on illegal token then print illegal syntax msg

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}
	//dir_name = filepath.Base(dir_path)
	//asm_file_name = dir_name + ".asm"
	//asm_path = dir_path + asm_file_name

	// create codewriter obj and send file to open to write
	//cw := codewriter.New(asm_path)

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".jack" {
			// Add actions here
			// loop through line by line of jack file and get list of tokens
			// then call compilaiton engine
			fmt.Println("End of input file: " + file.Name())
			// fmt.Println("Output file is ready: " + asm_file_name)
		}
	}
}
