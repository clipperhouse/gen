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
		err = get()
		// TODO: pass subsequent flags (such as -u) to get
	case "help":
		err = help()
	case "list":
		err = list()
	}

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
