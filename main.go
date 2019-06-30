package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
		position, _ := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(str))
		binaryPosition := fmt.Sprintf("%b", position)
		b := "0"
		for i := 1; i < 16 - len(binaryPosition); i++ {
			b = b + "0"
		}
		return b + binaryPosition
	}
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
