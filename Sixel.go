package main

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"strings"
)

const SixelBegin = "\x1bPq\n#0;2;0;0;0#1;2;100;100;100\n"
const SixelEnd = "\x1b\\"
const ANSI_WHITE = "\033[47m  \033[0m"
const ANSI_BLACK = "\033[40m  \033[0m"
const Block_Size = 6

func SixelPrint(output string) {
	println(SixelBegin)
	defer println(SixelEnd)
	for _, s := range strings.Split(output, "\n") {
		s = strings.ReplaceAll(s, ANSI_WHITE, fmt.Sprintf("#1!%d~", Block_Size))
		s = strings.ReplaceAll(s, ANSI_BLACK, fmt.Sprintf("#0!%d~", Block_Size))
		s = strings.ReplaceAll(s, " ", fmt.Sprintf("#0!%d~", Block_Size))
		for i := 0; i < Block_Size/6; i++ {
			print(s)
			println("-")
		}
	}
}

func TestSixelSupport(file *os.File) bool {
	//Send Control Character to Terminal and get response
	//If response contains DCS then Sixel is supported
	if !term.IsTerminal(int(file.Fd())) {
		return false

	}
	_, err := file.Write([]byte("\x1B[c"))
	if err != nil {
		return false
	}
	buf := make([]byte, 1024)
	//set echo off
	raw, err := term.MakeRaw(int(file.Fd()))
	defer term.Restore(int(file.Fd()), raw)
	_, err = file.Read(buf)
	if err != nil {
		return false
	}
	for _, b := range string(buf) {
		if b == '4' {
			//Found Sixel Support
			return true
		}
	}
	return false
}
