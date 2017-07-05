package parser

import (
	"strings"

	"github.com/eaburns/peggy/peg"
)

func walk(n *peg.Fail, f func(*peg.Fail) bool) bool {
	if !f(n) {
		return false
	}
	for _, k := range n.Kids {
		if !walk(k, f) {
			return false
		}
	}
	return true
}

// isWordNode returns whether the Node represents a word rule.
// Word rules are those whose names are all caps with no _, but is not solely the letter "h".
func isWordNode(n interface{}) bool {
	switch n := n.(type) {
	case *peg.Fail:
		return n.Name != "h" && isCaps(n.Name)
	case *peg.Node:
		return n.Name != "h" && isCaps(n.Name)
	default:
		panic("bad node type")
	}
}

const caps = "hABCDEFGIJKLMNOPRSTUVXYZ"

func isCaps(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(caps, r) {
			return false
		}
	}
	return len(s) > 0
}
