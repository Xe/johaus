package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"within.website/johaus/parser"
	"within.website/johaus/pretty"

	// Register all supported Lojban dialects in init().
	_ "within.website/johaus/parser/alldialects"
)

var (
	dialect        = flag.String("d", "camxes", "the dialect, one of: "+dialectString)
	keepMorph      = flag.Bool("m", false, "whether to keep morphology")
	addTerminators = flag.Bool("t", false, "whether to add elided terminators")
)

var dialectString = func() string {
	var s string
	for _, d := range parser.Dialects() {
		if s != "" {
			s += ", "
		}
		s += d.Name
	}
	return s
}()

func main() {
	flag.Parse()

	var r io.Reader
	var filePath string
	if len(flag.Args()) > 0 {
		filePath = flag.Arg(0)
		f, err := os.Open(filePath)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	text := string(data)
	fmt.Println("parsing")
	begin := time.Now()
	tree, err := parser.Parse(*dialect, text)
	end := time.Now()

	fmt.Println(end.Sub(begin))

	if err != nil {
		err.(*parser.Error).FilePath = filePath
		fmt.Println(err)
		os.Exit(1)
	}

	if !*keepMorph {
		parser.RemoveMorphology(tree)
	}
	if *addTerminators {
		parser.AddElidedTerminators(tree)
	}
	parser.RemoveSpace(tree)
	parser.CollapseLists(tree)

	pretty.Braces(os.Stdout, tree)
	fmt.Println("")

	pretty.Tree(os.Stdout, tree)
	fmt.Println("")
}
