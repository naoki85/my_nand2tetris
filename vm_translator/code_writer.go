package main

import (
	"os"
	"fmt"
	"strconv"
)

type CodeWriter struct {
	outputFileStream *os.File
}

func InitializeCodeWriter(fileName string) CodeWriter {
	c := CodeWriter{}
	c.SetFileName(fileName)
	return c
}

func (c *CodeWriter) SetFileName(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("err: Could not init parser %s", err.Error())
		os.Exit(1)
	}
	c.outputFileStream = file
}

func (c *CodeWriter) writeArithmetic(command string) {
	var output string
	if command == "add" {
		output = "@SP\nM=M-1\nA=M\n"
		output = output + "D=M\n"
		output = output + "@SP\nM=M-1\nA=M\n"
		output = output + "D=D+M\n"
		output = output + "@SP\n"
		output = output + "A=M\n"
		output = output + "M=D\n"
		output = output + "@SP\n"
		output = output + "M=M+1"
	}
	c.write(output)
}

func (c *CodeWriter) writePushPop(command string, segment string, index int) {
	var output string
	if command == CPush {
		output = "@" + strconv.Itoa(index) + "\n"
		output = output + "D=A\n"
		output = output + "@SP\n"
		output = output + "A=M\n"
		output = output + "M=D\n"
		output = output + "@SP\n"
		output = output + "M=M+1"
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
