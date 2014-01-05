package main

import (
	"errors"
	"fmt"
	"regexp"
)

type typeArg struct {
	Pointer, Package, Name string
}

func newTypeArg(s string) (result *typeArg, err error) {
	r := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)
	if matches := r.FindStringSubmatch(s); matches != nil {
		result = &typeArg{matches[1], matches[2], matches[3]}
		return
	}
	err = errors.New(fmt.Sprintf("could not parse %s into a type specification", s))
	return
}

func (t typeArg) String() string {
	return t.Pointer + t.Package + "." + t.Name
}

// "Package-local" name, includes pointer but ignores package.
func (t typeArg) LocalName() string {
	return t.Pointer + t.Name
}
