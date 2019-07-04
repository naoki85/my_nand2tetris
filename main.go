package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

	symbolTable, _ := initializeSymbolTable()
	setDefinedSymbol(symbolTable)

	for true {
		if regexp.MustCompile(`^@`).MatchString(parser.row) {
			label := strings.Replace(parser.row, "@", "", 1)
			if !symbolTable.contains(label) && !symbolTable.contains("R"+label) {
				symbolTable.addEntry(label, symbolTable.variableAddressCounter)
				symbolTable.variableAddressCounter += 1
			}
		}
		if !parser.hasMoreCommand() {
			break
		}
		parser.advance()
	}

	parser.ResetRowNumber()

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

func setDefinedSymbol(s SymbolTable) {
	s.addEntry("SP", 0)
	s.addEntry("LCL", 1)
	s.addEntry("ARG", 2)
	s.addEntry("THIS", 3)
	s.addEntry("THAT", 4)
	s.addEntry("SCREEN", 16384)
	s.addEntry("KBD", 24576)

	for i := 0; i < 16; i++ {
		label := "R" + strconv.Itoa(i)
		s.addEntry(label, i)
	}
}
