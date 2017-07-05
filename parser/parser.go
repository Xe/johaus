// Package parser contains the interface for using parsers for all supported Lojban dialects.
package parser

import (
	"errors"
	"fmt"
	"os"

	"github.com/eaburns/peggy/peg"
)

// Parser is a low-level interface to a Lojban parser.
type Parser interface {
	// Parse parses the text and returns
	// the maximum error position seen during the parse
	// and whether the parse succeeded.
	Parse() (int, bool)

	// ErrorTree returns the parse error tree for a failed parse.
	// The tree contains all errors at or beyond minErrorPos.
	ErrorTree(minErrorPos int) *peg.Fail

	// ParseTree returns the parse tree for a successful parse.
	ParseTree() *peg.Node
}

// Parse parses text using a given Lojban dialect.
// On success, the parseTree is returned.
// On failure, both the word-level and the raw, morphological errors are returned.
func Parse(dialect string, text string) (*peg.Node, error) {
	makeParser, ok := makeParserFuncs[dialect]
	if !ok {
		return nil, errors.New("unknown dialect: " + dialect)
	}
	p := makeParser(text)
	if perr, ok := p.Parse(); !ok {
		errTree := p.ErrorTree(perr)
		word := errorWordStart(errTree)
		if word < 0 {
			// BUG: {ji gi'e} with -d=maftufa returns a bad tree.
			fmt.Println("perr:", perr, "word err:", word)
			peg.PrettyWrite(os.Stdout, errTree)
			panic("")
		}
		errTree = p.ErrorTree(word)
		return nil, wordError(text, errTree)
	}
	return p.ParseTree(), nil
}
