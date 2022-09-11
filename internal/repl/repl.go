package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/vancanhuit/monkey/internal/evaluator"
	"github.com/vancanhuit/monkey/internal/lexer"
	"github.com/vancanhuit/monkey/internal/object"
	"github.com/vancanhuit/monkey/internal/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		prorgam := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		value := evaluator.Eval(prorgam, env)
		if value != nil {
			io.WriteString(out, value.Inspect())
			io.WriteString(out, "\n")
		}

		// io.WriteString(out, prorgam.String())
		// io.WriteString(out, "\n")
	}
}
