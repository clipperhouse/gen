package typewriter

import (
	"io"
)

type TypeWriter interface {
	Name() string
	// Validate is called for every Type to a) indicate that it will write for this Type and b) ensure that the TypeWriter considers it valid
	Validate(t Type) (bool, error)
	// WriteHeader writer to the top of the generated code, befoe the package declaration; intended for licenses.
	WriteHeader(w io.Writer, t Type)
	// Imports is a slice of import paths required for the type.
	Imports(t Type) []string
	// Write writes to the body of the generated code. This is the meat.
	WriteBody(w io.Writer, t Type)
}
