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

		if !known {
			err = errors.New("unknown argument: " + s)
			return
		}
	}

	return
}
