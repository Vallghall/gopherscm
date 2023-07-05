package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/Vallghall/gopherscm/internal/lexer"
	"github.com/Vallghall/gopherscm/internal/parser"
)

func main() {
	// going simple for now
	if len(os.Args) < 2 {
		log.Fatalln("Expected file name")
	}

	bs, err := os.ReadFile(os.Args[1])
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
	lexOut, _ := os.Create("lex.out.json")
	defer lexOut.Close()

	encoder := json.NewEncoder(lexOut)
	encoder.SetIndent("", "    ")
	encoder.Encode(ts)

	ast := parser.Parse(ts)

	parseOut, _ := os.Create("parse.out.json")
	defer parseOut.Close()

	encoder = json.NewEncoder(parseOut)
	encoder.SetIndent("", "    ")
	encoder.Encode(ast)

	_, err = ast.Eval()
	if err != nil {
		log.Fatalln(err)
	}
}
