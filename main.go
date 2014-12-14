// gen is a tool for type-driven code generation for Go. Details and docs are available at https://clipperhouse.github.io/gen.
package main

import (
	"os"
	"regexp"
)

func main() {
	var err error

	defer func() {
		if err != nil {
			if !exitStatusMsg.MatchString(err.Error()) {
				os.Stderr.WriteString(err.Error() + "\n")
			}
			os.Exit(1)
		}
	}()

	err = runMain(os.Args)
}

var exitStatusMsg = regexp.MustCompile(`^exit status \d+$`)

func runMain(args []string) error {
	c := defaultConfig

	if len(args) == 1 {
		// simply typed 'gen'; run is the default command
		return run(c)
	}

	cmd := args[1]

	var tail []string
	if len(args) > 2 {
		tail = args[2:]
	}

	switch cmd {
	case "add":
		return add(c, tail...)
	case "get":
		return get(c, tail...)
	case "list":
		return list(c)
	case "watch":
		return watch(c)
	default:
		return help(c)
	}
}
