package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("err: 1 argument is required")
		return
	}
	filePath := os.Args[1]
	outputFilePath := "test_files/Add.hack"

	parser, err := initializeParser(filePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		return
	}

	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		return
	}
	defer file.Close()

	for true {
		output := compileToBinary(parser)
		if len(output) != 16 {
			fmt.Printf("err: Could not parse code %s", output)
			break
		}
		_, err = file.Write(([]byte)(output + "\n"))
		if err != nil {
			fmt.Printf("err: Could not write code %s", err.Error())
			break
		}
		if !parser.hasMoreCommand() {
			break
		}
		parser.advance()
	}
}

func compileToBinary(p Parser) string {
	switch p.commandType() {
	case ACOMMAND:
		return p.getAddress()
	case CCOMMAND:
		var binary string
		code := Code{}
		binary = "111"
		binary = binary + code.comp(p.comp())
		binary = binary + code.dest(p.dest())
		binary = binary + code.jump(p.jump())
		return binary
	default:
		return ""
	}
}
