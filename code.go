package main

import "regexp"

type Code struct {
}

func (c Code) dest(mnemonic string) string {
	binary := ""
	if regexp.MustCompile(`A`).MatchString(mnemonic) {
		binary = binary + "1"
	} else {
		binary = binary + "0"
	}
	if regexp.MustCompile(`D`).MatchString(mnemonic) {
		binary = binary + "1"
	} else {
		binary = binary + "0"
	}
	if regexp.MustCompile(`M`).MatchString(mnemonic) {
		binary = binary + "1"
	} else {
		binary = binary + "0"
	}
	return binary
}

func (c Code) comp(mnemonic string) string {
	if regexp.MustCompile(`M`).MatchString(mnemonic) {
		return c.memoryComp(mnemonic)
	} else {
		return c.registerComp(mnemonic)
	}
}

func (c Code) jump(mnemonic string) string {
	switch mnemonic {
	case "null":
		return "000"
	case "JGT":
		return "001"
	case "JEQ":
		return "010"
	case "JGE":
		return "011"
	case "JLT":
		return "100"
	case "JNE":
		return "101"
	case "JLE":
		return "110"
	case "JMP":
		return "111"
	default:
		return ""
	}
}

func (c Code) memoryComp(mnemonic string) string {
	switch mnemonic {
	case "M":
		return "1110000"
	case "!M":
		return "1110001"
	case "-M":
		return "1110011"
	case "M+1":
		return "1110111"
	case "M-1":
		return "1110010"
	case "D+M":
		return "1000010"
	case "D-M":
		return "1010011"
	case "M-D":
		return "1000111"
	case "D&M":
		return "1000000"
	case "D|M":
		return "1010101"
	default:
		return ""
	}
}

func (c Code) registerComp(mnemonic string) string {
	switch mnemonic {
	case "0":
		return "0101010"
	case "1":
		return "0111111"
	case "-1":
		return "0111010"
	case "D":
		return "0001100"
	case "A":
		return "0110000"
	case "!D":
		return "0001101"
	case "!A":
		return "0110001"
	case "-D":
		return "0001111"
	case "-A":
		return "0110011"
	case "D+1":
		return "0011111"
	case "A+1":
		return "0110111"
	case "D-1":
		return "0001110"
	case "A-1":
		return "0110010"
	case "D+A":
		return "0000010"
	case "D-A":
		return "0010011"
	case "A-D":
		return "0000111"
	case "D&A":
		return "0000000"
	case "D|A":
		return "0010101"
	default:
		return ""
	}
}
