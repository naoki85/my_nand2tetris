package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("err: 1 argument is required")
		os.Exit(1)
	}
	filePath := os.Args[1]
	parser, err := InitializeParser(filePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		os.Exit(2)
	}

	outputFilePath := strings.Replace(filePath, ".vm", ".asm", 1)
	codeWriter := InitializeCodeWriter(outputFilePath)
	defer codeWriter.Close()

	for true {
		if parser.CommandType() == CPush {
			index, err := strconv.Atoi(parser.Arg2())
			if err != nil {
				fmt.Printf("err: Could not convert integer from string %s", err.Error())
				os.Exit(3)
			}
			codeWriter.writePushPop(parser.CommandType(), parser.Arg1(), index)
		} else if parser.CommandType() == CArithmetic {
			codeWriter.writeArithmetic(parser.row)
		}

		if !parser.HasMoreCommand() {
			break
		}
		parser.Advance()
	}
}
