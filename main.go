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
		output := compileToBinary(parser) + "\n"
		_, err = file.Write(([]byte)(output))
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
	default:
		str := p.row
		if str == "D=A" {
			return "1110110000010000"
		}
		if str == "D=D+A" {
			return "1110000010010000"
		}
		if str == "M=D" {
			return "1110001100001000"
		}
		return str
	}
}
