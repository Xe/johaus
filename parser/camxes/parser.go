// Package camxes is the camxes implementation of the official Lojban grammar.
package camxes

//go:generate peggy -o camxes.go camxes.peg

import (
	"github.com/eaburns/johaus/parser"
	"github.com/eaburns/peggy/peg"
)

func init() {
	parser.Register(
		parser.Dialect{
			Name:    "camxes",
			Version: "",
			Descr: map[string]string{
				"en":  "The camxes implementation of the official Lojban grammar.",
				"jbo": "la .camxes. cu ralju loi bankle be la .lojban.",
			},
			OfficialURL: "http://lojban.github.io/ilmentufa/camxes.html",
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
