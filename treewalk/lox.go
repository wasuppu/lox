package lox

import (
	"bufio"
	"fmt"
	"io"
	"lox/treewalk/astprinter"
	"lox/treewalk/interpreter"
	"lox/treewalk/loxerrors"
	"lox/treewalk/parser"
	"lox/treewalk/scanner"
	"os"
)

const PROMPT = ">> "

type lox struct {
	printer  *astprinter.ASTPrinter
	loxerror *loxerrors.LoxErrors
}

func New() *lox {
	return &lox{loxerror: loxerrors.New(), printer: astprinter.New()}
}

func (l *lox) RunPrompt(in io.Reader, out io.Writer) {
	l.loxerror = loxerrors.New()
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "" {
			break
		}

		l.run(line)
		l.loxerror.HadError = false
	}
}

func (l lox) RunFile(filename string) {
	l.loxerror = loxerrors.New()

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %q: %v", filename, err)
		os.Exit(1)
	}

	if len(data) > 0 {
		l.run(string(data))

	} else {
		fmt.Println("EOF null")
	}
}

func (l lox) run(source string) error {
	scanner := scanner.New(source, l.loxerror)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens, l.loxerror)
	interpreter := interpreter.New(l.loxerror)

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	if l.loxerror.HadError {
		os.Exit(65)
	}

	if l.loxerror.HadRuntimeError {
		os.Exit(70)
	}

	// l.printer.Print(exp)
	fmt.Println(l.printer.PrintStmts(statements))
	// interpreter.Interpret(statements)

	value, err := interpreter.Interpret(statements)
	if err != nil {
		return err
	}
	if value != nil {
		fmt.Println(astprinter.Stringify(value))
	}
	return nil
}
