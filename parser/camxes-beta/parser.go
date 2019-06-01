// Package camxes is the camxes implementation of the official Lojban grammar.
package camxes

//go:generate peggy -o camxes_beta.go camxes-beta.peg

import (
	"github.com/eaburns/peggy/peg"
	"within.website/johaus/parser"
)

func init() {
	parser.Register(
		parser.Dialect{
			Name:    "camxes-beta",
			Version: "",
			Descr: map[string]string{
				"en":  "The beta camxes implementation of the official Lojban grammar.",
				"jbo": "la .camxes. be lo cipra cu ralju loi bankle be la .lojban.",
			},
			OfficialURL: "http://lojban.github.io/ilmentufa/camxes.html",
			GrammarURL:  "",
		},
		func(text string) parser.Parser { return _NewParser(text) },
	)
}

func (p *_Parser) Parse() (int, bool) {
	pos, perr := _textAccepts(p, 0)
	return perr, pos >= 0
}

func (p *_Parser) ErrorTree(minPos int) *peg.Fail {
	p.fail = make(map[_key]*peg.Fail) // reset fail memo table
	_, tree := _textFail(p, 0, minPos)
	return tree
}

func (p *_Parser) ParseTree() *peg.Node {
	_, tree := _textNode(p, 0)
	return tree
}
