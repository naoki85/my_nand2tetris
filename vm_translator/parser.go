package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
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
