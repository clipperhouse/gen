package main

import (
	"regexp"
)

type options struct {
	All          bool
	AllPointer   string
	ExportedOnly bool
	Force        bool
}

func getOptions(args []string) *options {
	opts := &options{}

	allOption := regexp.MustCompile(`^-(\*?)a(ll)?$`)
	exportedOption := regexp.MustCompile(`^-e(xported)?$`)
	forceOption := regexp.MustCompile(`^-f(orce)?$`)

	for _, a := range args {
		allMatches := allOption.FindStringSubmatch(a)
		if allMatches != nil {
			opts.All = true
			opts.AllPointer = allMatches[1]
		}
		if exportedOption.MatchString(a) {
			opts.ExportedOnly = true
		}
		if forceOption.MatchString(a) {
			opts.Force = true
		}
	}

	return opts
}
