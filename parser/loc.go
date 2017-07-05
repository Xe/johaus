package parser

import (
	"sort"
	"unicode/utf8"

	"github.com/eaburns/peggy/peg"
)

// A Loc is a location in the input text.
type Loc struct {
	Byte   int
	Rune   int
	Line   int
	Column int
}

// Location returns the Loc at the corresponding byte offset in the text.
func Location(text string, byte int) Loc {
	var loc Loc
	loc.Line = 1
	loc.Column = 1
	for byte > loc.Byte {
		r, w := utf8.DecodeRuneInString(text[loc.Byte:])
		loc.Byte += w
		loc.Rune++
		loc.Column++
		if r == '\n' {
			loc.Line++
			loc.Column = 1
		}
	}
	return loc
}

// Locations returns a mapping from nodes to their Locs.
func Locations(text string, n *peg.Fail) map[*peg.Fail]Loc {
	var loc Loc
	loc.Line = 1
	loc.Column = 1
	locs := make(map[*peg.Fail]Loc)
	for _, n := range nodesSortedByPos(n) {
		for n.Pos > loc.Byte {
			r, w := utf8.DecodeRuneInString(text[loc.Byte:])
			loc.Byte += w
			loc.Rune++
			loc.Column++
			if r == '\n' {
				loc.Line++
				loc.Column = 1
			}
		}
		locs[n] = loc
	}
	return locs
}

func nodesSortedByPos(n *peg.Fail) []*peg.Fail {
	var nodes []*peg.Fail
	walk(n, func(n *peg.Fail) bool {
		nodes = append(nodes, n)
		return true
	})
	sort.Slice(nodes, func(i, j int) bool {
		a, b := nodes[i], nodes[j]
		return a.Pos == b.Pos && a.Name < b.Name || a.Pos < b.Pos
	})
	return nodes
}
