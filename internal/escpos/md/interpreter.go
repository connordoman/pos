package md

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var interpreter *Interpreter

type Interpreter struct {
	Verbose bool

	Parser *Parser
	Tokens []*Token

	hadError bool
}

func NewInterpreter() *Interpreter {
	if interpreter == nil {
		interpreter = &Interpreter{
			hadError: false,

			Parser: nil,
			Tokens: []*Token{},

			Verbose: false,
		}
	}
	return interpreter
}

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	log.Printf("[line %d] Error %s: %s", line, where, message)
	interpreter.hadError = true
}

func (i *Interpreter) RunFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	log.Println("Running file:", path)
	log.Printf("Length: %d bytes", len(data))

	err = i.Run(string(data))
	if err != nil {
		return err
	}

	if i.hadError {
		os.Exit(65)
		return fmt.Errorf("had error")
	}

	return nil
}

func (i *Interpreter) RunPrompt() error {
	var source string
	for {
		fmt.Print("→ ")
		if _, err := fmt.Scanln(&source); err != nil {
			return err
		}
		if strings.TrimSpace(source) == "" {
			continue
		}
		if err := i.Run(source); err != nil {
			fmt.Println("error:", err)
		}
	}
}

func (i *Interpreter) Run(source string) error {
	i.Parser = NewParser(source)
	i.Tokens = i.Parser.ScanTokens(source)

	for _, t := range i.Tokens {
		fmt.Println(t)
	}

	return nil
}
