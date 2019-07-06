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

	parser, err := InitializeParser(filePath)
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

	symbolTable, _ := InitializeSymbolTable()
	setDefinedSymbol(symbolTable)

	for true {
		if parser.CommandType() == LCommand {
			label := parser.Symbol()
			symbolTable.AddEntry(label, symbolTable.variableAddressCounter)
		} else {
			symbolTable.variableAddressCounter += 1
		}
		if !parser.HasMoreCommand() {
			break
		}
		parser.Advance()
	}

	parser.ResetRowNumber()

	for true {
		if parser.CommandType() == LCommand {
			parser.Advance()
			continue
		}
		output := compileToBinary(&parser, symbolTable)
		if len(output) != 16 {
			fmt.Printf("err: Could not parse code %s", parser.row)
			break
		}
		_, err = file.Write(([]byte)(output + "\n"))
		if err != nil {
			fmt.Printf("err: Could not write code %s", err.Error())
			break
		}
		if !parser.HasMoreCommand() {
			break
		}
		parser.Advance()
	}
}

func compileToBinary(p *Parser, s SymbolTable) string {
	switch p.CommandType() {
	case ACommand:
		return p.GetAddress(s)
	case CCommand:
		var binary string
		code := Code{}
		binary = "111"
		binary = binary + code.Comp(p.Comp())
		binary = binary + code.Dest(p.Dest())
		binary = binary + code.Jump(p.Jump())
		return binary
	default:
		return ""
	}
}

func setDefinedSymbol(s SymbolTable) {
	s.AddEntry("SP", 0)
	s.AddEntry("LCL", 1)
	s.AddEntry("ARG", 2)
	s.AddEntry("THIS", 3)
	s.AddEntry("THAT", 4)
	s.AddEntry("SCREEN", 16384)
	s.AddEntry("KBD", 24576)

	for i := 0; i < 16; i++ {
		label := "R" + strconv.Itoa(i)
		s.AddEntry(label, i)
	}
}
