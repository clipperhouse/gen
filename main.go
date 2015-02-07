// gen is a tool for type-driven code generation for Go. Details and docs are available at https://clipperhouse.github.io/gen.
package main

import (
	"fmt"
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

var exitStatusMsg = regexp.MustCompile("^exit status \\d+$")

func runMain(args []string) error {
	c := defaultConfig

	cmd, force, tail, err := parseArgs(args)

	if err != nil {
		return err
	}

	c.IgnoreTypeCheckErrors = force

	if len(cmd) == 0 {
		// simply typed 'gen'; run is the default command
		return run(c)
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

var s = struct{}{}

var cmds = map[string]struct{}{
	"add":   s,
	"get":   s,
	"help":  s,
	"list":  s,
	"watch": s,
}

func parseArgs(args []string) (cmd string, force bool, tail []string, err error) {
	for _, a := range args[1:] { // arg[0] is 'gen'
		if _, ok := cmds[a]; ok {
			if len(cmd) > 0 {
				err = fmt.Errorf("more than one command specified; type gen help for usage")
				break
			}
			cmd = a
			continue
		}
		if a == "-f" {
			force = true
			continue
		}
		tail = append(tail, a)
	}

	// tail is only valid on add & get; otherwise an error
	if len(tail) > 0 && cmd != "add" && cmd != "get" {
		err = fmt.Errorf("unknown command(s) %v", tail)
		tail = []string{}
	}

	// force flag is only valid with run & watch
	if force && cmd != "" && cmd != "watch" {
		err = fmt.Errorf("-f flag is not valid with %q", cmd)
	}

	return cmd, force, tail, err
}
