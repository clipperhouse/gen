package typewriter

import (
	"io"
)

type TypeWriter interface {
	Name() string
	// Validate is called for every Type to a) indicate that it will write for this Type and b) ensure that the TypeWriter considers it valid
	Validate(t Type) (bool, error)
	// Write to the top of the generated code; intended for license, or package-level comments.
	WriteHeader(w io.Writer, t Type)
	// Imports is a slice of names of imports required for the type.
	Imports(t Type) []string
	// Writes writes to the body of the generated code. This is the meat.
	Write(w io.Writer, t Type)
}
