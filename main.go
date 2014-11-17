// gen is a tool for type-driven code generation for Go. Details and docs are available at https://clipperhouse.github.io/gen.
package main

import "os"

func main() {
	var err error

	defer func() {
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}
	}()

	err = runMain(os.Args)
}

func runMain(args []string) error {
	if len(args) == 1 {
		// simply typed 'gen'; run is the default command
		return run()
	}

	cmd := args[1]

	switch cmd {
	case "custom":
		return custom()
	case "get":
		var tail []string
		if len(args) > 2 {
			tail = args[2:]
		}
		return get(tail)
	case "list":
		return list()
	default:
		return help()
	}
}
