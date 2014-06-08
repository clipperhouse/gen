package typewriter

import (
	"io"
	"os"
	"strings"
	"testing"
)

var fw *fooWriter = &fooWriter{}
var bw *barWriter = &barWriter{}

func TestRegister(t *testing.T) {
	// these registrations are stateful; remain in place for later tests
	// so keep 'em at the top of this file

	if err := Register(fw); err != nil {
		t.Error(err)
	}

	if err := Register(bw); err != nil {
		t.Error(err)
	}

	if err := Register(&fooWriter{}); err == nil {
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

	// see TestRegister; this really only tests that they've been assigned to the app
	if len(a1.TypeWriters) != 2 {
		t.Errorf("should have found 2 typewriters, found %v", len(a1.TypeWriters))
	}

	a2, err2 := NewApp("+dummy")

	if err2 != nil {
		t.Error(err2)
	}

	if len(a2.Types) != 0 {
		t.Errorf("should have found no types, found %v", len(a2.Types))
	}

	// see TestRegister; this really only tests that nothing has changed
	if len(a2.TypeWriters) != 2 {
		t.Errorf("should have found 2 typewriters, found %v", len(a2.TypeWriters))
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

func TestWriteAll(t *testing.T) {
	a1, err1 := NewApp("+test")

	if err1 != nil {
		t.Error(err1)
	}

	a1.WriteAll()

	if fw.validateCalls != len(a1.Types) {
		t.Errorf(".Validate() should have been called %v times (once for each type); was called %v", len(a1.Types), fw.validateCalls)
	}

	if bw.validateCalls != len(a1.Types) {
		t.Errorf(".Validate() should have been called %v times (once for each type); was called %v", len(a1.Types), bw.validateCalls)
	}

	if fw.writeHeaderCalls != len(a1.Types) {
		t.Errorf(".WriteHeader() should have been called %v times (once for each type); was called %v", len(a1.Types), fw.writeHeaderCalls)
	}

	// see Validate implementation below; chooses not to write dummy
	if bw.writeHeaderCalls != len(a1.Types)-1 {
		t.Errorf(".WriteHeader() should have been called %v times (once for each type); was called %v", len(a1.Types), bw.writeHeaderCalls)
	}

	if fw.writeCalls != len(a1.Types) {
		t.Errorf(".Write() should have been called %v times (once for each type); was called %v", len(a1.Types), fw.writeCalls)
	}

	// see Validate implementation below; chooses not to write dummy
	if bw.writeCalls != len(a1.Types)-1 {
		t.Errorf(".Write() should have been called %v times (once for each type); was called %v", len(a1.Types), bw.writeCalls)
	}
}

type fooWriter struct {
	validateCalls, writeHeaderCalls, writeCalls int
}

func (f *fooWriter) Name() string {
	return "foo"
}

func (f *fooWriter) Validate(t Type) (bool, error) {
	f.validateCalls++
	return true, nil
}

func (f *fooWriter) WriteHeader(w io.Writer, t Type) {
	f.writeHeaderCalls++
	return
}

func (f *fooWriter) Imports(t Type) (result []string) {
	return result
}

func (f *fooWriter) Write(w io.Writer, t Type) {
	f.writeCalls++
	return
}

type barWriter struct {
	validateCalls, writeHeaderCalls, writeCalls int
}

func (f *barWriter) Name() string {
	return "bar"
}

func (f *barWriter) Validate(t Type) (bool, error) {
	f.validateCalls++
	return t.Name != "dummy", nil // indicates not going to write for the dummy type
}

func (f *barWriter) WriteHeader(w io.Writer, t Type) {
	f.writeHeaderCalls++
	return
}

func (f *barWriter) Imports(t Type) (result []string) {
	return result
}

func (f *barWriter) Write(w io.Writer, t Type) {
	f.writeCalls++
	return
}
