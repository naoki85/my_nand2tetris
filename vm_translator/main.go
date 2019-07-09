package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("err: 1 argument is required")
		os.Exit(1)
	}
	filePath := os.Args[1]
	outputFilePath := strings.Replace(filePath, ".vm", ".asm", 1)

	parser, err := InitializeParser(filePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		os.Exit(2)
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		os.Exit(3)
	}
	defer file.Close()

	for true {
		var output string
		if regexp.MustCompile(`^push constant \d+`).MatchString(parser.row) {
			code := regexp.MustCompile(`\d+`).FindString(parser.row)
			output = "@" + code + "\n"
			output = output + "D=A\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		} else if parser.row == "add" {
			output = "@SP\nM=M-1\nA=M\n"
			output = output + "D=M\n"
			output = output + "@SP\nM=M-1\nA=M\n"
			output = output + "D=D+M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		}

		_, err = file.Write(([]byte)(output + "\n"))
		if err != nil {
			fmt.Printf("err: Could not write code %s", err.Error())
			os.Exit(5)
		}
		
		if !parser.HasMoreCommand() {
			break
		}
		parser.Advance()
	}
}
