package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("err: 1 argument is required")
		return
	}
	filepath := os.Args[1]
	outputFilePath := "test_files/Add.hack"

	fp, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("err: Could not open %s", filepath)
		return
	}
	defer fp.Close()

	file, err := os.Create(outputFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		if skipReadingRow(scanner.Text()) {
			continue
		}
		output := compileToBinary(scanner.Text()) + "\n"
		file.Write(([]byte)(output))
	}
	if err = scanner.Err(); err != nil {
		return
	}
}

func skipReadingRow(text string) bool {
	if len(text) == 0 || regexp.MustCompile(`^//`).MatchString(text) {
		return true
	}
	return false
}

func compileToBinary(str string) string {
	if regexp.MustCompile(`^@\d`).MatchString(str) {
		return "0000000000010000"
	}
	if str == "D=A" {
		return "1110110000001000"
	}
	if str == "D=D+A" {
		return "1110000010001000"
	}
	if str == "M=D" {
		return "1111110000010000"
	}
	return str
}
