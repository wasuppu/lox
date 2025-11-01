package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

var exprAnnotations = []string{
	"Literal : Value any",
	"Grouping : Expression Expr",
	"Unary : Operator token.Token, Right Expr",
	"Logical : Left Expr, Operator token.Token, Right Expr",
	"Binary : Left Expr, Operator token.Token, Right Expr",
	"Call : Callee Expr, Paren token.Token, Arguments []Expr",
	"Variable : Name token.Token",
	"Assign : Name token.Token, Value Expr",
}

var stmtAnnotations = []string{
	"Print : Expression Expr",
	"Return : Keyword token.Token, Value Expr",
	"Var : Name token.Token , Initializer Expr",
	"Block : Statements []Stmt",
	"Expression : Expression Expr",
	"Function : Name token.Token, Params []token.Token, Body []Stmt",
	"If : Condition Expr, ThenBranch Stmt, ElseBranch Stmt",
	"While : Condition Expr, Body Stmt",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: astgen <output>")
		os.Exit(1)
	}

	out, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	defineAst(out, "Expr", exprAnnotations)
	defineAst(out, "Stmt", stmtAnnotations)
}

func defineAst(out, base string, types []string) {
	source := "package ast\n"
	source += `import "glox/treewalk/token"`
	source += fmt.Sprintln()

	source += defineVisitor(base, types)

	source += fmt.Sprintf(`
type %s interface {
	Accept(v %sVisitor) any
}
    `, base, base)
	source += fmt.Sprintln()

	for _, t := range types {
		name := strings.TrimRight(strings.Split(t, ":")[0], " ")
		fields := strings.Trim(strings.Split(t, ":")[1], " ")
		source += defineType(base, name, fields)
	}

	path := fmt.Sprintf("%s/%s.go", out, strings.ToLower(base))
	if err := saveFile(path, source); err != nil {
		panic(err)
	}
	// fmt.Println(source)
}

func defineVisitor(base string, types []string) string {
	source := fmt.Sprintf("type %sVisitor interface {\n", base)

	for _, t := range types {
		name := strings.TrimRight(strings.Split(t, ":")[0], " ")
		source += fmt.Sprintf("Visit%s%s(expr *%s) any\n", name, base, name)
	}
	source += "}\n"
	return source
}

func defineType(base, name, fields string) string {
	var source string

	source += fmt.Sprintf("type %s struct {\n", name)

	// fields
	flds := strings.Split(fields, ",")
	source += strings.Join(flds, "\n")
	source += fmt.Sprintln("\n}")

	// New func
	rets := []string{}
	args := []string{}

	for _, fld := range flds {
		t := strings.Split(strings.Trim(fld, " "), " ")[1]
		s := strings.Split(strings.Trim(fld, " "), " ")[0]
		e := strings.ToLower(s)

		args = append(args, fmt.Sprintf("%s %s", e, t))
		rets = append(rets, fmt.Sprintf("%s: %s", s, e))
	}
	argument := fmt.Sprint(strings.Join(args, ","))
	source += fmt.Sprintf("func New%s(%s) %s {\n", name, argument, base)

	mems := fmt.Sprint(strings.Join(rets, ","))
	source += fmt.Sprintf("return &%s{%s}\n}\n", name, mems)

	// Accept func
	source += fmt.Sprintf(`
func (e *%s) Accept(v %sVisitor) any {
    return v.Visit%s%s(e)
}
`, name, base, name, base)

	return source
}

func saveFile(path, source string) error {
	// gofmt
	buf, err := format.Source([]byte(source))
	if err != nil {
		return err
	}

	os.WriteFile(path, buf, 0644)

	return err
}
