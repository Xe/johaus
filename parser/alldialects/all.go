// Package alldialects can be imported to register all supported Lojban dialects.
package alldialects

// Import all supported parsers for registration.
import (
	_ "within.website/johaus/parser/camxes"
	_ "within.website/johaus/parser/ilmentufa"
	_ "within.website/johaus/parser/maftufa"
	_ "within.website/johaus/parser/zantufa"
)
