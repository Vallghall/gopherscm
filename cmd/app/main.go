package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Vallghall/gopherscm/internal/lexer"
	"github.com/Vallghall/gopherscm/internal/parser"
)

func main() {
	ops := getOptions()

	// going simple for now
	if len(flag.Args()) < 1 {
		log.Fatalln("Expected file name")
	}

	bs, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalln("** FILE NOT FOUND **")
		}

		log.Fatalln(err)
	}

	ts, err := lexer.Lex(bytes.Runes(bs))
	if err != nil {
		log.Fatalln(err)
	}

	if *ops.lexOut {
		writeIntermediateResults("lex", ts)
	}

	ast := parser.Parse(ts)

	if *ops.parseOut {
		writeIntermediateResults("parse", ast)
	}

	_, err = ast.Eval()
	if err != nil {
		log.Fatalln(err)
	}
}

// options - wrapper around flag options
type options struct {
	lexOut   *bool
	parseOut *bool
}

// getOptions - helper for retrieving flag values
func getOptions() *options {
	defer flag.Parse()
	return &options{
		lexOut:   flag.Bool("L", false, "logs lexer's results into a lex.out.json file"),
		parseOut: flag.Bool("P", false, "logs parser's results into a parse.out.json file"),
	}
}

// writeIntermediateResults - helper for writing intermediate results,
// such as lexer's token stream or parser's ats into a json file
func writeIntermediateResults(prefix string, result any) {
	out, err := os.Create(fmt.Sprintf("%s.out.json", prefix))
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()

	encoder := json.NewEncoder(out)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(result)
	if err != nil {
		log.Println(err)
	}
}
