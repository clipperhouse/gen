package main

import (
	"os"

	"github.com/codegangsta/cli"
)

const customFilename string = "_gen.go"

func main() {
	a := &cli.App{
		Name:    os.Args[0],
		Usage:   "http://clipperhouse.github.io/gen",
		Version: "3.0.0",
		Author:  "Matt Sherman",
		Email:   "mwsherman@gmail.com",
		Action: func(c *cli.Context) {
			run()
		},
		Commands: []cli.Command{
			// keep UI (cli) concerns out of the main routines
			{
				Name: "custom",
				Action: func(c *cli.Context) {
					custom(customFilename)
				},
				Usage: "Creates a custom _gen.go file in which to specify your own typewriter imports",
			},

			{
				Name: "get",
				Action: func(c *cli.Context) {
					get(c.Bool("u"))
				},
				Usage: "Runs `go get` for gen typewriters; intended for custom typewriters in _gen.go; unnecessary when using the defaults",
				Flags: []cli.Flag{
					cli.BoolFlag{"u", "use the network to update the typewriter packages and their dependencies"},
				},
			},
		},
	}

	a.Run(os.Args)
}
