// Package alldialects can be imported to register all supported Lojban dialects.
package alldialects

// Import all supported parsers for registration.
import (
	_ "github.com/Xe/johaus/parser/camxes"
	_ "github.com/Xe/johaus/parser/ilmentufa"
	_ "github.com/Xe/johaus/parser/maftufa"
	_ "github.com/Xe/johaus/parser/zantufa"
)
