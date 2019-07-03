package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Parser struct {
	rows      []string
	row       string
	rowNumber int
}

const (
	ACOMMAND string = "A_COMMAND"
	CCOMMAND string = "C_COMMAND"
	LCOMMAND string = "L_COMMAND"
)

func initializeParser(filePath string) (Parser, error) {
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
		p.rows = append(p.rows, scanner.Text())
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

func (p *Parser) hasMoreCommand() bool {
	if p.rowNumber+1 < len(p.rows) {
		return true
	}
	return false
}

func (p *Parser) advance() {
	p.rowNumber = p.rowNumber + 1
	p.row = p.rows[p.rowNumber]
}

func (p *Parser) commandType() string {
	if regexp.MustCompile(`^@\d`).MatchString(p.row) {
		return ACOMMAND
	}
	if regexp.MustCompile(`\s*(=|;)\s*`).MatchString(p.row) {
		return CCOMMAND
	}
	switch p.row {
	case LCOMMAND:
		return LCOMMAND
	default:
		return ""
	}
}

func (p *Parser) symbol() {
	log.Fatal("Must be implemented")
}

func (p *Parser) dest() string {
	slice := regexp.MustCompile(`\s*=\s*`).Split(p.row, 2)
	if len(slice) != 2 {
		return "null"
	}
	return slice[0]
}

func (p *Parser) comp() string {
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

func (p *Parser) jump() string {
	sliceMnemonic := regexp.MustCompile(`\s*;\s*`).Split(p.row, 2)
	if len(sliceMnemonic) < 2 || sliceMnemonic[0] == "" {
		return "null"
	}
	return sliceMnemonic[1]
}

func (p *Parser) getAddress() string {
	position, _ := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(p.row))
	binaryPosition := fmt.Sprintf("%b", position)
	b := "0"
	for i := 1; i < 16-len(binaryPosition); i++ {
		b = b + "0"
	}
	return b + binaryPosition
}
