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
		return
	}
	filePath := os.Args[1]

	parser, err := initializeParser(filePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		return
	}

	inputFileName := regexp.MustCompile(`[0-9a-zA-Z_]*.asm$`).FindString(filePath)
	outputFilePath := strings.Split(inputFileName, ".")[0] + ".hack"

	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		return
	}
	defer file.Close()

	for true {
		output := compileToBinary(parser)
		if len(output) != 16 {
			fmt.Printf("err: Could not parse code %s", parser.row)
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
