package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/eaburns/peggy/peg"
	"within.website/johaus/parser"
	_ "within.website/johaus/parser/alldialects"
	"within.website/johaus/pretty"
)

func init() {
	http.HandleFunc("/", rootHandler)
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	dialect, err := parserDialect(req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseGlob("*.tmplt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := map[string]interface{}{
			"Dialect":  dialect,
			"Dialects": dialectNames,
		}
		if err := t.ExecuteTemplate(w, "parser.tmplt", data); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		text, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		resp := make(map[string]interface{})
		tree, err := parser.Parse(dialect.Name, string(text))
		if err != nil {
			resp["Error"] = err.Error()
		} else {
			query := req.URL.Query()
			if q := query["morph"]; len(q) < 1 || q[0] != "true" {
				parser.RemoveMorphology(tree)
			}
			if q := query["terms"]; len(q) > 0 && q[0] == "true" {
				parser.AddElidedTerminators(tree)
			}
			parser.RemoveSpace(tree)
			parser.CollapseLists(tree)
			resp["Tree"] = prettyString(pretty.Tree, tree)
			resp["Braces"] = prettyString(pretty.Braces, tree)
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	default:
		msg := "method " + req.Method + " is not allowed"
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}

}

func prettyString(printer func(io.Writer, *peg.Node) error, tree *peg.Node) string {
	buf := bytes.NewBuffer(nil)
	printer(buf, tree)
	return buf.String()
}

// parserDialect looks up the parser.Dialect for the requested parser.
// If the dialect is not supported, a user-readable error is returned.
func parserDialect(url *url.URL) (*parser.Dialect, error) {
	parserName := "camxes"
	if b := path.Base(url.Path); b != "/" {
		parserName = b
	}
	for _, d := range parser.Dialects() {
		if d.Name == parserName {
			return &d, nil
		}
	}
	return nil, errors.New(parserName + " is not supported. Supported dialects are: " + strings.Join(dialectNames, ", "))
}

var dialectNames = func() []string {
	var ns []string
	for _, d := range parser.Dialects() {
		ns = append(ns, d.Name)
	}
	return ns
}()
