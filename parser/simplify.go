package parser

import (
	"strings"

	"github.com/eaburns/peggy/peg"
)

// RemoveSpace removes whitespace-only nodes.
func RemoveSpace(n *peg.Node) { removeSpace(n) }

func removeSpace(n *peg.Node) bool {
	if whitespace(n.Text) {
		return false
	}
	if len(n.Kids) == 0 {
		return true
	}
	var kids []*peg.Node
	for _, k := range n.Kids {
		if removeSpace(k) {
			kids = append(kids, k)
		}
	}
	n.Kids = kids
	return len(n.Kids) > 0
}

// SpaceChars is the string of all whitespace characters.
const SpaceChars = ".\t\n\r?!\x20"

func whitespace(s string) bool {
	for _, r := range s {
		if !strings.ContainsRune(SpaceChars, r) {
			return false
		}
	}
	return true
}

// CollapseLists collapses chains of single-kid nodes.
func CollapseLists(n *peg.Node) {
	if collapseLists(n) == 1 {
		n.Kids = n.Kids[0].Kids
	}
}

func collapseLists(n *peg.Node) int {
	var kids []*peg.Node
	for _, k := range n.Kids {
		if gk := collapseLists(k); gk == 1 {
			kids = append(kids, k.Kids[0])
		} else {
			kids = append(kids, k)
		}
	}
	n.Kids = kids
	return len(n.Kids)
}

// AddElidedTerminators sets the Text of an elided terminator node to the terminator name in all caps.
// This flags any functions removing empty Nodes to keep the elided terminator Node, as its Text is no longer empty.
func AddElidedTerminators(n *peg.Node) {
	const elidableSuffix = "_elidible" // elidable is spelled wrong in the PEG grammar files.

	if n.Text == "" && strings.HasSuffix(n.Name, elidableSuffix) {
		n.Text = strings.TrimSuffix(n.Name, elidableSuffix)
		n.Kids = nil
		return
	}
	for _, k := range n.Kids {
		AddElidedTerminators(k)
	}
}

// RemoveMorphology removes all nodes beneath whole words.
func RemoveMorphology(n *peg.Node) {
	if isWordNode(n) || n.Name == "zoi_word" {
		n.Kids = nil
		return
	}
	for _, k := range n.Kids {
		RemoveMorphology(k)
	}
}
