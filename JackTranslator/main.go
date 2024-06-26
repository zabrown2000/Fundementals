package main

import (
	"Fundementals/JackTranslator/compilationEngine"
	"Fundementals/JackTranslator/tokeniser"
	"fmt"
	"os"
	"path"
	"strings"
)

var CurJACK string

func main() {
	// get path from user
	fmt.Println("Enter path to folder")
	var dir_path string
	_, err := fmt.Scanln(&dir_path)
	if err != nil {
		return
	}

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

			hierarchOutPath := strings.TrimSuffix(dir_path+CurJACK, ".jack") + "New.xml"
			plainOutPath := strings.TrimSuffix(dir_path+CurJACK, ".jack") + "NewT.xml"

			ce := compilationEngine.New(plainOutPath, hierarchOutPath, tkn)
			ce.CompileClass()
			fmt.Println("End of input file: " + file.Name())
		}
	}
}
