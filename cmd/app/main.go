package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	ast := parser.Parse(ts)
	j, _ := json.MarshalIndent(ast, "", "    ")
	fmt.Println(string(j))
}
