// add catch for if parser throws failure on illegal token then print illegal syntax msg
package main

import (
	"Fundementals/JackTranslator/tokeniser"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var CurJACK string
var asm_path string
var dir_name string
var asm_file_name string

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}
	dir_name = filepath.Base(dir_path)

	files, err := os.ReadDir(dir_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if path.Ext(file.Name()) == ".jack" {
			CurJACK = file.Name()
			// each vm file, create parser obj and send file to open to read
			tkn, err := tokeniser.New(dir_path + CurJACK)
			if err != nil {
				fmt.Println(err)
				return
			}
			tkn.TokeniseLine()
			fmt.Println("End of input file: " + file.Name())
		}
	}
}
