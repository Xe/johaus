// Package zantufa is the zantufa dialect of Lojban.
package zantufa

//go:generate peggy -o zantufa.go zantufa-1.9999.peg

import (
	"github.com/eaburns/peggy/peg"
	"within.website/johaus/parser"
)

func init() {
	parser.Register(
		parser.Dialect{
			Name:    "zantufa",
			Version: "1.9999",
			Descr: map[string]string{
				"en":  "The zantufa dialect of Lojban.",
				"jbo": "la zantufa cu ciplanli bankle la .lojban.",
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
