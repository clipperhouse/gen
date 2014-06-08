package typewriter

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	if err := Register(fooWriter{}); err != nil {
		t.Error(err)
	}

	if err := Register(barWriter("baz")); err != nil {
		t.Error(err)
	}

	if err := Register(fooWriter{}); err == nil {
		t.Error("registering the same typewriter twice should be an error")
	}

	if len(typeWriters) != 2 {
		t.Error("should have 2 typewriters registered, found %v", len(typeWriters))
	}
}

func TestNewApp(t *testing.T) {
	a1, err1 := NewApp("+test")

	if err1 != nil {
		t.Error(err1)
	}

	// app and dummy types
	if len(a1.Types) != 2 {
		t.Errorf("should have found 2 types, found", len(a1.Types))
	}

	a2, err2 := NewApp("+dummy")

	if err2 != nil {
		t.Error(err2)
	}

	if len(a2.Types) != 0 {
		t.Errorf("should have found no types, found", len(a2.Types))
	}
}

func TestNewAppFiltered(t *testing.T) {
	filter := func(f os.FileInfo) bool {
		return !strings.HasPrefix(f.Name(), "app")
	}

	a1, err1 := NewAppFiltered("+test", filter)

	if err1 != nil {
		t.Error(err1)
	}

	// app is filtered out, only dummy type
	if len(a1.Types) != 1 {
		t.Errorf("should have found 1 types, found %v", len(a1.Types))
	}
}

type fooWriter struct{}

func (f fooWriter) Name() string {
	return "foo"
}

func (f fooWriter) Validate(t Type) (bool, error) {
	return true, nil
}

func (f fooWriter) WriteHeader(w io.Writer, t Type) {
	return
}

func (f fooWriter) Imports(t Type) (result []string) {
	return result
}

func (f fooWriter) Write(w io.Writer, t Type) {
	return
}

type barWriter string // heck, could be anything

func (f barWriter) Name() string {
	return "bar"
}

func (f barWriter) Validate(t Type) (bool, error) {
	return true, nil
}

func (f barWriter) WriteHeader(w io.Writer, t Type) {
	return
}

func (f barWriter) Imports(t Type) (result []string) {
	return result
}

func (f barWriter) Write(w io.Writer, t Type) {
	return
}
