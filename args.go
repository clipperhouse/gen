package main

import (
	"errors"
	"regexp"
)

type typeArg struct {
	Pointer, Package, Name string
}

func (t typeArg) String() string {
	return t.Pointer + t.Package + "." + t.Name
}

// "Package-local" name, includes pointer but ignores package.
func (t typeArg) LocalName() string {
	return t.Pointer + t.Name
}

type options struct {
	All          bool
	AllPointer   string
	ExportedOnly bool
	Force        bool
}

func parseArgs(args []string) (typeArgs []*typeArg, opts options, errs []error) {
	opts = options{}

	typ := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)
	all := regexp.MustCompile(`^-(\*?)a(ll)?$`)
	exported := regexp.MustCompile(`^-e(xported)?$`)
	force := regexp.MustCompile(`^-f(orce)?$`)

	for _, s := range args {
		known := false

		if matches := typ.FindStringSubmatch(s); matches != nil {
			t := &typeArg{matches[1], matches[2], matches[3]}
			typeArgs = append(typeArgs, t)
			known = true
		}

		if matches := all.FindStringSubmatch(s); matches != nil {
			opts.All = true
			opts.AllPointer = matches[1]
			known = true
		}

		if exported.MatchString(s) {
			opts.ExportedOnly = true
			known = true
		}

		if force.MatchString(s) {
			opts.Force = true
			known = true
		}

		if !known {
			errs = append(errs, errors.New("unknown argument: "+s))
		}
	}

	if opts.ExportedOnly && !opts.All {
		errs = append(errs, errors.New("the -e(xported) flag is only valid when used with the -a(ll) flag"))
	}

	if len(typeArgs) == 0 && !opts.All {
		errs = append(errs, errors.New("at least one type, or the -all flag, is required"))
	}

	if len(typeArgs) > 0 && opts.All {
		errs = append(errs, errors.New("either specify a type or use the -all flag; do not do both"))
	}

	return
}
