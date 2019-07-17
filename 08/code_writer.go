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
	var output string
	var symbolOutput string
	switch command {
	case "add", "sub":
		output = "@SP\nM=M-1\nA=M\nD=M\n"
		output = output + "@SP\nM=M-1\nA=M\n"
		if command == "add" {
			output = output + "D=D+M\n"
		} else {
			output = output + "D=M-D\n"
		}
		output = output + "@SP\n"
		output = output + "A=M\n"
		output = output + "M=D\n"
		output = output + "@SP\nM=M+1"
	case "neg":
		output = "@SP\n"
		output = output + "A=M-1\n"
		output = output + "M=-M"
	case "and":
		output = "@SP\nM=M-1\nA=M\nD=M\n"
		output = output + "@SP\nM=M-1\nA=M\n"
		output = output + "D=D&M\n"
		output = output + "@SP\n"
		output = output + "A=M\n"
		output = output + "M=D\n"
		output = output + "@SP\nM=M+1"
	case "or":
		output = "@SP\nM=M-1\nA=M\nD=M\n"
		output = output + "@SP\nM=M-1\nA=M\n"
		output = output + "D=D|M\n"
		output = output + "@SP\n"
		output = output + "A=M\n"
		output = output + "M=D\n"
		output = output + "@SP\nM=M+1"
	case "not":
		output = "@SP\n"
		output = output + "A=M-1\n"
		output = output + "M=!M"
	case "eq", "lt", "gt":
		output = "@SP\nM=M-1\nA=M\nD=M\n"
		output = output + "@SP\nM=M-1\nA=M\nD=M-D\n"
		numberStr := strconv.Itoa(c.labelNumber)
		output = output + "@LABEL" + numberStr + "\n"
		symbolOutput = "(LABEL" + numberStr + ")\nD=-1\n"
		c.labelNumber++

		if command == "eq" {
			output = output + "D;JEQ\nD=0\n"
		} else if command == "lt" {
			output = output + "D;JLT\nD=0\n"
		} else {
			output = output + "D;JGT\nD=0\n"
		}
		
		numberStr = strconv.Itoa(c.labelNumber)
		output = output + "@LABEL" + numberStr + "\n"
		symbolOutput = symbolOutput + "(LABEL" + numberStr + ")\n@SP\nA=M\nM=D\n"
		c.labelNumber++

		output = output + "0;JMP\n"
		output = output + symbolOutput
		output = output + "@SP\nM=M+1"
	default: output = ""
	}
	c.write(output + "\n")
}

func (c *CodeWriter) WritePushPop(command string, segment string, index int) {
	var output string
	if command == CPush {
		switch segment {
		case "constant":
			output = "@" + strconv.Itoa(index) + "\n"
			output = output + "D=A\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "local":
			output = output + "@LCL\n"
			output = output + "A=M\n"
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "that":
			output = output + "@THAT\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "this":
			output = output + "@THIS\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "pointer":
			output = output + "@3\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "argument":
			output = output + "@ARG\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		case "temp":
			output = output + "@5\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "D=M\n"
			output = output + "@SP\n"
			output = output + "A=M\n"
			output = output + "M=D\n"
			output = output + "@SP\n"
			output = output + "M=M+1"
		default:
			output = segment
		}
	} else if command == CPop {
		output = "@SP\n"
		output = output + "M=M-1\n"
		output = output + "A=M\n"
		output = output + "D=M\n"
		switch segment {
		case "local":
			output = output + "@LCL\n"
			output = output + "A=M\n"
			output = output + "M=D"
		case "argument":
			output = output + "@ARG\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "M=D"
		case "this":
			output = output + "@THIS\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "M=D"
		case "that":
			output = output + "@THAT\n"
			output = output + "A=M\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "M=D"
		case "pointer":
			output = output + "@3\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "M=D"
		case "temp":
			output = output + "@5\n"
			for i := 1; i <= index; i++ {
				output = output + "A=A+1\n"
			}
			output = output + "M=D"
		default:
		}
	}
	c.write(output + "\n")
}

func (c *CodeWriter) Close() {
	c.outputFileStream.Close()
}

func (c *CodeWriter) write(output string) {
	_, err := c.outputFileStream.Write(([]byte)(output))
	if err != nil {
		fmt.Printf("err: Could not write code %s", err.Error())
		os.Exit(5)
	}
}

func (c *CodeWriter) pushFromDRegister() {
	slice := []string{
		"@SP", "A=M", "M=D", "@SP", "M=M+1",
	}
	c.writeCodes(slice)
}

func (c *CodeWriter) writeCodes(slice []string) {
	output := strings.Join(slice, "\n")
	_, err := c.outputFileStream.Write(([]byte)(output + "\n"))
	if err != nil {
		fmt.Printf("err: Could not write code %s", err.Error())
		os.Exit(5)
	}
}
