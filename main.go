package main

import "os"

func main() {
	args := os.Args

	if len(args) == 1 {
		// simply typed 'gen'; run is the default command
		runCmd(run)
		return
		// TODO: exclude test files?
	}

	cmd := args[1]

	switch cmd {
	case "custom":
		runCmd(custom)
	case "get":
		runCmd(get)
		// TODO: pass subsequent flags (such as -u) to get
	case "list":
		runCmd(list)
	}
	// TODO: add help command
	// TODO: verbosity?
}

func runCmd(f func() error) {
	err := f()

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
