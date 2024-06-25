// add catch for if parser throws failure on illegal token then print illegal syntax msg
package main

import (
	"Fundementals/JackTranslator/compilationEngine"
	"Fundementals/JackTranslator/tokeniser"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
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
			tkn.TokeniseFile()
			//fmt.Println(tkn.LengthTokens)
			//for i := 0; i < tkn.LengthTokens-1; i++ {
			//	//fmt.Println(tkn.Tokens[i].Token_type)
			//	//fmt.Println(tkn.Tokens[i].Token_content)
			//}

			hierarchOutPath := strings.TrimSuffix(dir_path+CurJACK, ".jack") + "New.xml"
			plainOutPath := strings.TrimSuffix(dir_path+CurJACK, ".jack") + "NewT.xml"
			//fileOut, err := os.Create(hierarchOutPath)
			//if err != nil {
			//	fmt.Println("Failed to create hierarch output file:", hierarchOutPath)
			//	return
			//}
			//defer fileOut.Close()
			//
			//tokenFileOut, err := os.Create(plainOutPath)
			//if err != nil {
			//	fmt.Println("Failed to create plain output file:", plainOutPath)
			//	return
			//}
			//defer tokenFileOut.Close()
			ce := compilationEngine.New(plainOutPath, hierarchOutPath, tkn)
			ce.CompileClass()
			fmt.Println("End of input file: " + file.Name())
		}
	}
}
