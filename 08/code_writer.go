package main

import (
	"os"
	"fmt"
	"strconv"
	"strings"
)

type CodeWriter struct {
	outputFileStream *os.File
	labelNumber int
	labelReturnNumber int
}

func InitializeCodeWriter(fileName string) CodeWriter {
	c := CodeWriter{}
	c.SetFileName(fileName)
	c.labelNumber = 1
	c.labelReturnNumber = 0
	c.init()
	return c
}

func (c *CodeWriter) init() {
	c.writeCodes([]string{
		fmt.Sprintf("@%d", 256),
		"D=A",
		"@SP",
		"M=D",
	})

	// return label
	c.labelReturnNumber++
	label := "_RETURN_LABEL_" + strconv.Itoa(c.labelReturnNumber)
	c.writeCodes([]string{fmt.Sprintf("@%s", label), "D=A"})

	c.pushFromDRegister()
	c.writeCodes([]string{"@LCL", "D=M"})
	c.pushFromDRegister()
	c.writeCodes([]string{"@ARG", "D=M"})
	c.pushFromDRegister()
	c.writeCodes([]string{"@THIS", "D=M"})
	c.pushFromDRegister()
	c.writeCodes([]string{"@THAT", "D=M"})
	c.pushFromDRegister()
	c.writeCodes([]string{
		"@SP",
		"D=M",
		"@5",
		"D=D-A",
		"@0",
		"D=D-A",
		"@ARG",
		"M=D",
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
	})

	c.writeCodes([]string{
		"@Sys.init",
		"0;JMP",
		fmt.Sprintf("(%s)", label),
	})
}

func (c *CodeWriter) SetFileName(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		os.Exit(1)
	}
	c.outputFileStream = file
}

func (c *CodeWriter) WriteArithmetic(command string) {
	switch command {
	case "add", "sub":
		var output string
		c.popToMRegister()
		c.writeCodes([]string{"D=M"})
		c.popToMRegister()
		if command == "add" {
			output = "D=D+M"
		} else {
			output = "D=M-D"
		}
		c.writeCodes([]string{output})
		c.pushFromDRegister()
	case "neg":
		c.writeCodes([]string{"@SP", "A=M-1", "M=-M"})
	case "and":
		c.popToMRegister()
		c.writeCodes([]string{"D=M"})
		c.popToMRegister()
		c.writeCodes([]string{"D=D&M"})
		c.pushFromDRegister()
	case "or":
		c.popToMRegister()
		c.writeCodes([]string{"D=M"})
		c.popToMRegister()
		c.writeCodes([]string{"D=D|M"})
		c.pushFromDRegister()
	case "not":
		c.writeCodes([]string{"@SP", "A=M-1", "M=!M"})
	case "eq", "lt", "gt":
		c.popToMRegister()
		c.writeCodes([]string{"D=M"})
		c.popToMRegister()
		c.writeCodes([]string{"D=M-D"})
		var outputSlice []string
		var symbolOutput []string
		numberStr := strconv.Itoa(c.labelNumber)
		outputSlice = append(outputSlice, "@LABEL" + numberStr)
		symbolOutput = append(symbolOutput, "(LABEL" + numberStr + ")")
		symbolOutput = append(symbolOutput, "D=-1")
		c.labelNumber++

		if command == "eq" {
			outputSlice = append(outputSlice, "D;JEQ")
		} else if command == "lt" {
			outputSlice = append(outputSlice, "D;JLT")
		} else {
			outputSlice = append(outputSlice, "D;JGT")
		}
		outputSlice = append(outputSlice, "D=0")
		numberStr = strconv.Itoa(c.labelNumber)
		outputSlice = append(outputSlice, "@LABEL" + numberStr)
		symbolOutput = append(symbolOutput, "(LABEL" + numberStr + ")")
		c.labelNumber++

		outputSlice = append(outputSlice, "D;JMP")
		outputSlice = append(outputSlice, symbolOutput...)
		c.writeCodes(outputSlice)
		c.pushFromDRegister()
	default:
		fmt.Printf("warning: Invalid command %s", command)
	}
}

func (c *CodeWriter) WritePushPop(command string, segment string, index int) {
	if command == CPush {
		switch segment {
		case "constant":
			address := "@" + strconv.Itoa(index)
			c.writeCodes([]string{address, "D=A"})
			c.pushFromDRegister()
		case "local":
			c.writeCodes([]string{"@LCL", "A=M"})
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		case "that":
			c.writeCodes([]string{"@THAT", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		case "this":
			c.writeCodes([]string{"@THIS", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		case "pointer":
			c.writeCodes([]string{"@3"})
			c.increaseAddress(index)
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		case "argument":
			c.writeCodes([]string{"@ARG", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		case "temp":
			c.writeCodes([]string{"@5"})
			c.increaseAddress(index)
			c.writeCodes([]string{"D=M"})
			c.pushFromDRegister()
		default:
		}
	} else if command == CPop {
		c.writeCodes([]string{"@SP", "M=M-1", "A=M", "D=M"})
		switch segment {
		case "local":
			c.writeCodes([]string{"@LCL", "A=M", "M=D"})
		case "argument":
			c.writeCodes([]string{"@ARG", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"M=D"})
		case "this":
			c.writeCodes([]string{"@THIS", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"M=D"})
		case "that":
			c.writeCodes([]string{"@THAT", "A=M"})
			c.increaseAddress(index)
			c.writeCodes([]string{"M=D"})
		case "pointer":
			c.writeCodes([]string{"@3"})
			c.increaseAddress(index)
			c.writeCodes([]string{"M=D"})
		case "temp":
			c.writeCodes([]string{"@5"})
			c.increaseAddress(index)
			c.writeCodes([]string{"M=D"})
		default:
		}
	}
}

func (c *CodeWriter) Close() {
	c.outputFileStream.Close()
}

func (c *CodeWriter) pushFromDRegister() {
	slice := []string{
		"@SP", "A=M", "M=D", "@SP", "M=M+1",
	}
	c.writeCodes(slice)
}

func (c *CodeWriter) popToMRegister() {
	slice := []string{
		"@SP", "M=M-1", "A=M",
	}
	c.writeCodes(slice)
}

func (c *CodeWriter) increaseAddress(index int) {
	var output []string
	for i := 1; i <= index; i++ {
		output = append(output, "A=A+1")
	}
	c.writeCodes(output)
}

func (c *CodeWriter) writeCodes(slice []string) {
	output := strings.Join(slice, "\n")
	_, err := c.outputFileStream.Write(([]byte)(output + "\n"))
	if err != nil {
		fmt.Printf("err: Could not write code %s", err.Error())
		os.Exit(5)
	}
}
