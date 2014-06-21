package typewriter

import (
	"io"
)

type TypeWriter interface {
	Name() string
	// Validate is called for every Type, prior to further action, to answer two questions:
	// a) that your TypeWriter will write for this Type; return false if your TypeWriter intends not to write a file for this Type at all.
	// b) that your TypeWriter considers the declaration (i.e., Tags) valid; return err if not.
	Validate(t Type) (bool, error)
	// WriteHeader writer to the top of the generated code, before the package declaration; intended for licenses or general documentation.
	WriteHeader(w io.Writer, t Type)
	// Imports is a slice of imports required for the type; each will be written into the imports declaration.
	Imports(t Type) []ImportSpec
	// Write writes to the body of the generated code, following package declaration, headers and imports. This is the meat.
	WriteBody(w io.Writer, t Type)
}
