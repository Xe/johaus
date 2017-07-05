package parser

import "sort"

// A Dialect describes a dialect of the Lojban language.
type Dialect struct {
	// Name is the name of the dialect.
	Name string

	// Version is the version of the dialect's grammar.
	Version string

	// Descr is a map from language codes
	// to text describing the dialect
	// in the corresponding language.
	Descr map[string]string

	// OfficialURL is the official URL for the parser's homepage.
	OfficialURL string

	// GrammarURL is the URL of the grammar file.
	GrammarURL string
}

var (
	dialects        = make(map[string]Dialect)
	makeParserFuncs = make(map[string]func(string) Parser)
)

// Register registers a dialect.
func Register(d Dialect, makeParser func(string) Parser) {
	dialects[d.Name] = d
	makeParserFuncs[d.Name] = makeParser
}

// Dialects returns all registered dialects in lexical order by name.
func Dialects() []Dialect {
	var ds []Dialect
	for _, d := range dialects {
		ds = append(ds, d)
	}
	sort.Slice(ds, func(i, j int) bool {
		return ds[i].Name < ds[j].Name
	})
	return ds
}
