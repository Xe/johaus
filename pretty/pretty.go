// Package pretty provides functions for pretty-printing parse trees.
package pretty

import (
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/eaburns/peggy/peg"
)

var braces = [][2]string{
	{"(", ")"},
	{"[", "]"},
	{"{", "}"},
	{"<", ">"},
}

// Braces pretty-prints a parse tree using ( ), [ ], { }, and < > show the nesting structure.
func Braces(w io.Writer, n *peg.Node) error {
	var walk func(int, *peg.Node) error
	walk = func(depth int, n *peg.Node) error {
		if n.Text == "" {
			return nil
		}
		if len(n.Kids) == 0 {
			_, err := io.WriteString(w, n.Text)
			return err
		}
		if _, err := io.WriteString(w, braces[depth][0]); err != nil {
			return err
		}
		for i, kid := range n.Kids {
			walk((depth+1)%len(braces), kid)
			if i != len(n.Kids)-1 {
				if _, err := io.WriteString(w, " "); err != nil {
					return err
				}
			}
		}
		_, err := io.WriteString(w, braces[depth][1])
		return err
	}
	return walk(0, n)
}

// Tree writes a pretty representation of the tree.
func Tree(w io.Writer, n *peg.Node) error {
	nr := utf8.RuneCountInString(n.Name)
	return tree(w, "", nr, n)
}

func tree(w io.Writer, tab string, nameWidth int, n *peg.Node) error {
	if len(n.Kids) == 0 {
		_, err := io.WriteString(w, n.Name+"["+strconv.Quote(n.Text)+"]\n")
		return err
	}
	name := n.Name
	if name == "" {
		name = strings.Repeat("─", nameWidth/2) + "┐"
	}
	if _, err := io.WriteString(w, name+"\n"); err != nil {
		return err
	}
	kidNameWidth := -1
	for _, kid := range n.Kids {
		n := utf8.RuneCountInString(kid.Name)
		if n > 0 && (kidNameWidth < 0 || n < kidNameWidth) {
			kidNameWidth = n
		}
	}
	if kidNameWidth < 4 {
		kidNameWidth = 4
	}
	spaces := strings.Repeat(" ", nameWidth/2)
	line := strings.Repeat("─", nameWidth/2)
	for i, kid := range n.Kids {
		kidTab := tab + spaces
		kidLine := kidTab
		if i == len(n.Kids)-1 {
			kidTab += " " + spaces
			kidLine += "└" + line
		} else {
			kidTab += "│" + spaces
			kidLine += "├" + line
		}
		if _, err := io.WriteString(w, kidLine); err != nil {
			return err
		}
		if err := tree(w, kidTab, kidNameWidth, kid); err != nil {
			return err
		}
	}
	return nil
}
