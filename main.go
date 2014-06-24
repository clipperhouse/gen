package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

const customFilename string = "_gen.go"

func main() {
	// keep UI (cli) concerns out of the main routines
	// output and exit should happen up here, not down there
	// we send in customFilename and receive 'out' to be a bit more testable; otherwise would use a const and fmt.Println
	a := &cli.App{
		Name:    os.Args[0],
		Usage:   "http://clipperhouse.github.io/gen",
		Version: "3.0.0",
		Author:  "Matt Sherman",
		Email:   "mwsherman@gmail.com",
		Action: func(c *cli.Context) {
			err := run(customFilename)

			if err != nil {
				log.Fatalln(err)
			}
		},
		Commands: []cli.Command{
			{
				Name: "custom",
				Action: func(c *cli.Context) {
					if err := custom(customFilename); err != nil {
						log.Fatalln(err)
					}
				},
				Usage: "Creates a custom _gen.go file in which to specify your own typewriter imports",
			},
			{
				Name: "get",
				Action: func(c *cli.Context) {
					if err := get(c.Bool("u")); err != nil {
						log.Fatalln(err)
					}
				},
				Usage: "Runs `go get` for gen typewriters; intended for custom typewriters in _gen.go; unnecessary when using the defaults",
				Flags: []cli.Flag{
					cli.BoolFlag{"u", "use the network to update the typewriter packages and their dependencies"},
				},
			},
			{
				Name: "list",
				Action: func(c *cli.Context) {
					err := list(customFilename)

					if err != nil {
						log.Fatalln(err)
					}
				},
				Usage: "Lists current typewriters",
			},
		},
	}

	a.Run(os.Args)
}
