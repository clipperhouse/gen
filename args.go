package main

import (
	"errors"
	"regexp"
)

type options struct {
	Force bool
	Help  bool
}

func parseArgs(args []string) (opts options, err error) {
	opts = options{}
	force := regexp.MustCompile(`^-f(orce)?$`)
	help := regexp.MustCompile(`^-h(elp)?$`)

	for _, s := range args {
		known := false

		if force.MatchString(s) {
			opts.Force = true
			known = true
		}

		if help.MatchString(s) {
			opts.Help = true
			known = true
		}

		// give informative errors for use of deprecated api
		// TODO: update URLs below after merge with master
		typ := regexp.MustCompile(`^(\*?)([\p{L}\p{N}]+)\.([\p{L}\p{N}]+)$`)
		if typ.MatchString(s) {
			err = errors.New("passing type arguments at the command line has been deprecated, see https://github.com/clipperhouse/gen/blob/projection/CHANGELOG.md")
			return
		}

		all := regexp.MustCompile(`^-(\*?)a(ll)?$`)
		if all.MatchString(s) {
			err = errors.New("the -all flag has been deprecated, see https://github.com/clipperhouse/gen/blob/projection/CHANGELOG.md")
			return
		}

		exported := regexp.MustCompile(`^-e(xported)?$`)
		if exported.MatchString(s) {
			err = errors.New("the -exported flag has been deprecated, see https://github.com/clipperhouse/gen/blob/projection/CHANGELOG.md")
			return
		}

		if !known {
			err = errors.New("unknown argument: " + s)
			return
		}
	}

	return
}
