// Package maftufa is the maftufa dialect of Lojban grammar.
package maftufa

//go:generate peggy -o maftufa.go maftufa.peg

import (
	"github.com/eaburns/johaus/parser"
	"github.com/eaburns/peggy/peg"
)

func init() {
	parser.Register(
		parser.Dialect{
			Name:    "maftufa",
			Version: "1.1",
			Descr: map[string]string{
				"en":  "The maftufa dialect of Lojban grammar, created for parsing lo se manci te makfa pe la oz (http://selpahi.de/oz_plain.html).",
				"jbo": "la .maftufa. cu ciplanli bankle la .lojban. te zu'e lo nu gentufa la'e lo se manci te makfa pe la .oz.",
			},
			OfficialURL: "https://mw.lojban.org/papri/zantufa",
			GrammarURL:  "",
		},
		func(text string) parser.Parser { return _NewParser(text) },
	)
}

func (p *_Parser) Parse() (int, bool) {
	pos, perr := _text_eofAccepts(p, 0)
	return perr, pos >= 0
}

func (p *_Parser) ErrorTree(minPos int) *peg.Fail {
	p.fail = make(map[_key]*peg.Fail) // reset fail memo table
	_, tree := _text_eofFail(p, 0, minPos)
	return tree
}

func (p *_Parser) ParseTree() *peg.Node {
	_, tree := _text_eofNode(p, 0)
	return tree
}
