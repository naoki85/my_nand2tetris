package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	CArithmetic = "C_ARITHMETIC"
	CPush = "C_PUSH"
	CPop = "C_POP"
	CLabel = "C_LABEL"
	CGoto = "C_GOTO"
	CIf = "C_IF"
	CFunction = "C_FUNCTION"
	CReturn = "C_Return"
	CCall = "C_CALL"
)

type Parser struct {
	rows       []string
	row        string
	rowNumber  int
}

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
	if (regexp.MustCompile(`^push`).MatchString(p.row)) {
		return CPush
	} else if (regexp.MustCompile(`^pop`).MatchString(p.row)) {
		return CPop
	} else if (regexp.MustCompile(`^label`).MatchString(p.row)) {
		return CLabel
	} else if (regexp.MustCompile(`^if-`).MatchString(p.row)) {
		return CIf
	} else {
		return CArithmetic
	}
}

func (p *Parser) Arg1() string {
	return strings.Split(p.row, " ")[1]
}

func (p *Parser) Arg2() string {
	if (p.CommandType() == CPush || p.CommandType() == CPop ||
	p.CommandType() == CFunction || p.CommandType() == CCall) {
		return strings.Split(p.row, " ")[2]
	} else {
		fmt.Println("err: Invalid command type")
		os.Exit(1)
	}
	return ""
}
