// Package ilmentufa is the ilmentufa dialect of Lojban.
package ilmentufa

//go:generate peggy -o ilmentufa.go ilmentufa.peg

import (
	"github.com/Xe/johaus/parser"
	"github.com/eaburns/peggy/peg"
)

func init() {
	parser.Register(
		parser.Dialect{
			Name:    "ilmentufa",
			Version: "",
			Descr: map[string]string{
				"en":  "The ilmentufa dialect of Lojban.",
				"jbo": "la ilmentufa cu ciplanli bankle la .lojban.",
			},
			OfficialURL: "https://lojban.github.io/ilmentufa/glosser/glosser.htm",
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
