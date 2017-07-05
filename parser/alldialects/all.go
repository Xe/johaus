// Package alldialects can be imported to register all supported Lojban dialects.
package alldialects

// Import all supported parsers for registration.
import (
	_ "github.com/eaburns/johaus/parser/camxes"
	_ "github.com/eaburns/johaus/parser/ilmentufa"
	_ "github.com/eaburns/johaus/parser/maftufa"
	_ "github.com/eaburns/johaus/parser/zantufa"
)
