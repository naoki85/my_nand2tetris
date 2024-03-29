package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	rows       []string
	row        string
	rowNumber  int
	ramAddress int
}

const (
	ACommand string = "A_COMMAND"
	CCommand string = "C_COMMAND"
	LCommand string = "L_COMMAND"
)

func InitializeParser(filePath string) (Parser, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("err: Could not open %s", filePath)
		return Parser{}, err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	p := Parser{}
	for scanner.Scan() {
		if skipReadingRow(scanner.Text()) {
			continue
		}
		text := strings.Split(scanner.Text(), "//")[0]
		text = strings.TrimSpace(text)
		p.rows = append(p.rows, text)
	}
	if err = scanner.Err(); err != nil {
		return p, err
	}
	p.rowNumber = 0
	p.row = p.rows[0]
	p.ramAddress = 16
	return p, nil
}

func skipReadingRow(text string) bool {
	if len(text) == 0 || regexp.MustCompile(`^//`).MatchString(text) {
		return true
	}
	return false
}

func (p *Parser) HasMoreCommand() bool {
	if p.rowNumber+1 < len(p.rows) {
		return true
	}
	return false
}

func (p *Parser) Advance() {
	p.rowNumber = p.rowNumber + 1
	p.row = p.rows[p.rowNumber]
}

func (p *Parser) CommandType() string {
	if regexp.MustCompile(`^@`).MatchString(p.row) {
		return ACommand
	}
	if regexp.MustCompile(`\s*(=|;)\s*`).MatchString(p.row) {
		return CCommand
	}
	if regexp.MustCompile(`^\(.*\)$`).MatchString(p.row) {
		return LCommand
	}
	return ""
}

func (p *Parser) Symbol() string {
	symbol := strings.Replace(p.row, "(", "", 1)
	symbol = strings.Replace(symbol, ")", "", 1)
	return symbol
}

func (p *Parser) Dest() string {
	slice := regexp.MustCompile(`\s*=\s*`).Split(p.row, 2)
	if len(slice) != 2 {
		return "null"
	}
	return slice[0]
}

func (p *Parser) Comp() string {
	if regexp.MustCompile(`\s*=\s*`).MatchString(p.row) {
		sliceMnemonic := regexp.MustCompile(`\s*=\s*`).Split(p.row, 2)
		if len(sliceMnemonic) < 2 || sliceMnemonic[1] == "" {
			return "null"
		}
		return sliceMnemonic[1]
	} else if regexp.MustCompile(`\s*;\s*`).MatchString(p.row) {
		sliceMnemonic := regexp.MustCompile(`\s*;\s*`).Split(p.row, 2)
		if len(sliceMnemonic) < 2 || sliceMnemonic[0] == "" {
			return "null"
		}
		return sliceMnemonic[0]
	}
	return "null"
}

func (p *Parser) Jump() string {
	sliceMnemonic := regexp.MustCompile(`\s*;\s*`).Split(p.row, 2)
	if len(sliceMnemonic) < 2 || sliceMnemonic[0] == "" {
		return "null"
	}
	return sliceMnemonic[1]
}

func (p *Parser) GetAddress(symbolTable SymbolTable) string {
	var position string
	position = strings.Replace(p.row, "@", "", 1)

	var address int
	var err error
	if symbolTable.Contains(position) {
		address = symbolTable.GetAddress(position)
	} else {
		address, err = strconv.Atoi(position)
		if err != nil {
			symbolTable.AddEntry(position, p.ramAddress)
			address = p.ramAddress
			p.ramAddress = p.ramAddress + 1
		}
	}
	binaryPosition := fmt.Sprintf("%b", address)
	b := "0"
	for i := 1; i < 16-len(binaryPosition); i++ {
		b = b + "0"
	}
	return b + binaryPosition
}

func (p *Parser) ResetRowNumber() {
	p.rowNumber = 0
	p.row = p.rows[p.rowNumber]
}
