package main

import "os"

func main() {
	var err error

	args := os.Args

	if len(args) == 1 {
		// simply typed 'gen'; run is the default command
		err = run()
		return
	}

	cmd := args[1]

	switch cmd {
	case "custom":
		err = custom()
	case "get":
		var tail []string
		if len(args) > 2 {
			tail = args[2:]
		}
		err = get(tail)
	case "list":
		err = list()
	case "help":
		fallthrough
	default:
		err = help()
	}

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
