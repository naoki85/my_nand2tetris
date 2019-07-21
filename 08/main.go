package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Print("err: 1 argument is required")
		os.Exit(1)
	}
	filePath := os.Args[1]
	if regexp.MustCompile(`.vm$`).MatchString(filePath) {
		parser, err := InitializeParser(filePath)
		if err != nil {
			fmt.Printf("err: Could not init parser %s", err.Error())
			os.Exit(2)
		}
	
		outputFilePath := strings.Replace(filePath, ".vm", ".asm", 1)
		codeWriter := InitializeCodeWriter(outputFilePath)
		defer codeWriter.Close()
		translate(parser, &codeWriter)
	} else if fInfo, _ := os.Stat(filePath); fInfo.IsDir() {
		files, _ := ioutil.ReadDir(filePath)
		sliceFilePath := strings.Split(filePath, "/")
		baseName := sliceFilePath[len(sliceFilePath) - 1]
		codeWriter := InitializeCodeWriter(filePath + "/" + baseName + ".asm")
		defer codeWriter.Close()
		for _, f := range files {
			if !regexp.MustCompile(`.vm$`).MatchString(f.Name()) { continue }
			parser, err := InitializeParser(filePath + "/" + f.Name())
			if err != nil {
				fmt.Printf("err: Could not init parser %s", err.Error())
				os.Exit(2)
			}
			translate(parser, &codeWriter)
		}
	} else {
		fmt.Println("Err: Invalid file path")
		os.Exit(2)
	}
}

func translate(parser Parser, codeWriter *CodeWriter) {
	for true {
		if parser.CommandType() == CPush || parser.CommandType() == CPop {
			index, err := strconv.Atoi(parser.Arg2())
			if err != nil {
				fmt.Printf("err: Could not convert integer from string %s", err.Error())
				os.Exit(3)
			}
			codeWriter.WritePushPop(parser.CommandType(), parser.Arg1(), index)
		} else if parser.CommandType() == CArithmetic {
			codeWriter.WriteArithmetic(parser.row)
		} else if parser.CommandType() == CIf {
			codeWriter.WriteIf(parser.row)
		} else if parser.CommandType() == CLabel {
			codeWriter.WriteLabel(parser.row)
		} else if parser.CommandType() == CGoto {
			codeWriter.WriteGoto(parser.row)
		} else if parser.CommandType() == CFunction {
			arg, _ := strconv.Atoi(parser.Arg2())
			codeWriter.WriteFunction(parser.Arg1(), arg)
		} else if parser.CommandType() == CReturn {
			codeWriter.WriteReturn()
		} else if parser.CommandType() == CCall {
			arg, _ := strconv.Atoi(parser.Arg2())
			codeWriter.WriteCall(parser.Arg1(), arg)
		}

		if !parser.HasMoreCommand() {
			break
		}
		parser.Advance()
	}
}
