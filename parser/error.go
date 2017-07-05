package parser

import (
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/eaburns/peggy/peg"
)

// A Error represents an error at some location in the input text.
type Error struct {
	Loc
	FilePath string
	Want     []string
}

func (err Error) Error() string {
	var want string
	for i, w := range err.Want {
		if len(want) > 0 {
			want += ", "
		}
		if len(want) > 1 && i == len(err.Want)-1 {
			want += "or "
		}
		want += w
	}
	if len(err.Want) > 1 {
		want = "one of: " + want
	} else {
		want = ": " + want
	}
	return fmt.Sprintf("%s:%d.%d: expected %s", err.FilePath, err.Line, err.Column, want)
}

// rawError returns an Error from a failed parse tree with the raw, morphological errors.
// The FilePath on the returned Error is the empty string,
// but can be set by the caller.
func rawError(text string, n *peg.Fail) *Error {
	fails := getLeaves(n)
	sort.Slice(fails, func(i, j int) bool {
		switch a, b := fails[i], fails[j]; {
		case namedError(a) && !namedError(b):
			return false
		case !namedError(a) && namedError(b):
			return true
		default:
			return wantString(a) < wantString(b)
		}
	})
	max := -1
	var wants []string
	for _, f := range fails {
		switch {
		case f.Pos > max:
			max = f.Pos
			wants = []string{wantString(f)}
		case f.Pos == max:
			wants = append(wants, wantString(f))
		}
	}
	return &Error{Loc: Location(text, max), Want: wants}
}

func getLeaves(n *peg.Fail) []*peg.Fail {
	uniq := make(map[*peg.Fail]bool)
	walk(n, func(n *peg.Fail) bool {
		if len(n.Kids) == 0 {
			if r, _ := utf8.DecodeRuneInString(n.Want); r != '&' && r != '!' {
				uniq[n] = true
			}
		}
		return true
	})
	leaves := make([]*peg.Fail, 0, len(uniq))
	for n := range uniq {
		leaves = append(leaves, n)
	}
	return leaves
}

// wordError returns a word-level Error from a failed parse tree.
// The FilePath on the returned Error is the empty string,
// but can be set by the caller.
func wordError(text string, n *peg.Fail) *Error {
	fails, _ := getFails(n)
	sort.Slice(fails, func(i, j int) bool {
		switch a, b := fails[i], fails[j]; {
		case namedError(a) && !namedError(b):
			return false
		case !namedError(a) && namedError(b):
			return true
		default:
			return wantString(a) < wantString(b)
		}
	})
	var wants []string
	pos := errorWordStart(n)
	seen := make(map[*peg.Fail]bool)
	for _, f := range fails {
		if f.Pos == pos && !seen[f] {
			seen[f] = true
			wants = append(wants, wantString(f))
		}
	}
	return &Error{Loc: Location(text, pos), Want: wants}
}

func wantString(n *peg.Fail) string {
	if name := prettyName(n); !namedError(n) && name != "" {
		return name
	}
	return n.Want
}

// namedError returns whether the fail node represents a named error in the grammar,
// a rule whose name is followed by an error name before the <-.
func namedError(n *peg.Fail) bool {
	if n.Want == "" {
		return false
	}
	// Peggy error node texts for non-named errors begin with one of: " . [ ! or &.
	r, _ := utf8.DecodeRuneInString(n.Want)
	return !strings.ContainsRune(`".[!&`, r)
}

func getFails(n *peg.Fail) ([]*peg.Fail, int) {
	if len(n.Kids) == 0 {
		return []*peg.Fail{n}, n.Pos
	}
	max := n.Pos
	var fails []*peg.Fail
	for _, k := range n.Kids {
		fs, m := getFails(k)
		fails = append(fails, fs...)
		if m > max {
			max = m
		}
	}
	pretty := prettyName(n)
	if isWordNode(n) || pretty != "" && max == n.Pos {
		return []*peg.Fail{n}, n.Pos
	}
	return fails, max
}

// prettyName returns the user-displayable name of the node if it has one.
// Nodes with pretty names are any that are not "text" and do not contain _, except for bridi_tail and sumti_tail, which have the pretty names "bridi tail" and "sumti tail" respectively.
// The pretty name of BRIVLA, CMEVLA, and CMAVO are brivla, cmevla, and cmavo respectively.
func prettyName(n *peg.Fail) string {
	if n.Name == "bridi_tail" || n.Name == "sumti_tail" {
		return strings.Replace(n.Name, "_", " ", 1)
	}
	if n.Name == "text" || strings.ContainsRune(n.Name, '_') {
		return ""
	}
	return n.Name
}

// ErrorWordStart returns the maximum byte offset into the input
// among all word-node errors.
func errorWordStart(n *peg.Fail) int {
	max := -1
	if isWordNode(n) {
		max = n.Pos
	}
	for _, k := range n.Kids {
		if m := errorWordStart(k); m > max {
			max = m
		}
	}
	return max
}
