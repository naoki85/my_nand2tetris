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
	currentFunctionName string
	currentFileName string
}

func InitializeCodeWriter(fileName string) CodeWriter {
	c := CodeWriter{}
	c.SetFileName(fileName)
	c.labelNumber = 1
	c.labelReturnNumber = 0
	c.currentFunctionName = ""
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

	c.WriteCall("Sys.init", 0)
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

		outputSlice = append(outputSlice, "0;JMP")
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
			c.increaseAddress(index)
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
		case "static":
			c.writeCodes([]string{
				fmt.Sprintf("@%s.%d", c.currentFileName, index),
				"D=M", "@SP", "A=M", "M=D", "@SP", "M=M+1",
			})
		default:
		}
	} else if command == CPop {
		c.popToMRegister()
		c.writeCodes([]string{"D=M"})
		switch segment {
		case "local":
			c.writeCodes([]string{"@LCL", "A=M"})
			c.increaseAddress(index)
		case "argument":
			c.writeCodes([]string{"@ARG", "A=M"})
			c.increaseAddress(index)
		case "this":
			c.writeCodes([]string{"@THIS", "A=M"})
			c.increaseAddress(index)
		case "that":
			c.writeCodes([]string{"@THAT", "A=M"})
			c.increaseAddress(index)
		case "pointer":
			c.writeCodes([]string{"@3"})
			c.increaseAddress(index)
		case "temp":
			c.writeCodes([]string{"@5"})
			c.increaseAddress(index)
		case "static":
			c.writeCodes([]string{
				fmt.Sprintf("@%s.%d", c.currentFileName, index),
			})
		}
		c.writeCodes([]string{"M=D"})
	}
}

func (c *CodeWriter) SetCurrentFileName(fileName string) {
	c.currentFileName = fileName
}

func (c *CodeWriter) WriteLabel(label string) {
	c.writeCodes([]string{
		fmt.Sprintf("(%s)", c.parseLabelLine(label)),
	})
}

func (c *CodeWriter) WriteIf(label string) {
	c.popToMRegister()
	c.writeCodes([]string{
		"D=M", fmt.Sprintf("@%s", c.parseLabelLine(label)), "D;JNE",
	})
}

func (c *CodeWriter) WriteGoto(label string) {
	c.writeCodes([]string{
		fmt.Sprintf("@%s", c.parseLabelLine(label)), "0;JMP",
	})
}

func (c *CodeWriter) WriteFunction(functionName string, numLocals int) {
	c.writeCodes([]string{
		fmt.Sprintf("(%s)", functionName), "D=0",
	})
	for i := 1; i <= numLocals; i++ {
		c.pushFromDRegister()
	}
	c.currentFunctionName = functionName
}

func (c *CodeWriter) WriteCall(functionName string, numLocals int) {
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
		fmt.Sprintf("@%d", numLocals),
		"D=D-A",
		"@ARG",
		"M=D",
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
	})

	c.writeCodes([]string{
		fmt.Sprintf("@%s", functionName),
		"0;JMP",
		fmt.Sprintf("(%s)", label),
	})
}

func (c *CodeWriter) WriteReturn() {
	c.writeCodes([]string{
		"@LCL", "D=M",
		// R13 = FRAME = LCL
		"@R13", "M=D",
		"@5", "D=A",
		// D = *(FRAME-5) = return-address
		"@R13", "A=M-D", "D=M",
		// R14 = return-address
		"@R14", "M=D",
	})
	c.popToMRegister()
	c.writeCodes([]string{
		"D=M", "@ARG",
		// M = *ARG
		"A=M",
		// *ARG = pop()
		"M=D",
		"@ARG", "D=M+1",
		// SP = ARG + 1
		"@SP", "M=D",
		// A = FRAME-1, R13 = FRAME-1
		"@R13", "AM=M-1", "D=M",
		// THAT = *(FRAME-1)
		"@THAT", "M=D",
		"@R13", "AM=M-1", "D=M",
		// THIS = *(FRAME-2)
		"@THIS", "M=D",
		"@R13", "AM=M-1", "D=M",
		// ARG = *(FRAME-3)
		"@ARG", "M=D",
		"@R13", "AM=M-1", "D=M",
		// LCL = *(FRAME-4)
		"@LCL", "M=D",
		// goto return-address
		"@R14", "A=M", "0;JMP",
	})
}

func (c *CodeWriter) WriteComment(line string) {
	c.writeCodes([]string{fmt.Sprintf("// %s", line)})
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
	if index < 1 { return }
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

func (c *CodeWriter) parseLabelLine(label string) string {
	var functionName string
	if c.currentFunctionName == "" {
		functionName = "null"
	} else {
		functionName = c.currentFunctionName
	}
	text := strings.Split(label, " ")
	return fmt.Sprintf("%s$%s", functionName, text[1])
}
