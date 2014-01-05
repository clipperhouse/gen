package main

import (
	"strings"
)

// Utility for handling various string representations of types. Does no validation, ensure that you initialize/convert with something like *package.Type. Pointer and package are optional.
type typeString string

const ptr = "*"

func (t typeString) Pointer() string {
	if strings.HasPrefix(string(t), ptr) {
		return ptr
	}
	return ""
}

func (t typeString) Package() string {
	parts := strings.Split(string(t), ".")
	if len(parts) > 1 {
		s := parts[0]
		if strings.HasPrefix(s, ptr) {
			return s[1:]
		}
		return s
	}
	return ""
}

// name of the type only, no pointer or package
func (t typeString) Name() string {
	s := string(t)
	parts := strings.Split(s, ".")
	if len(parts) > 1 {
		return parts[1]
	}
	if strings.HasPrefix(s, ptr) {
		return s[1:]
	}
	return s
}

// "Package-local" name, includes pointer but ignores package.
func (t typeString) LocalName() string {
	return t.Pointer() + t.Name()
}
