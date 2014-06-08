package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// set up basic command args
	run := []string{
		"run",
	}

	// set up temp directory
	caller := filepath.Base(os.Args[0])

	tempDir, err := ioutil.TempDir("", caller)

	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tempDir)

	// read gen custom imports file
	if src, err := ioutil.ReadFile("_gen.go"); err == nil {
		custom := filepath.Join(tempDir, "custom.go")

		// write custom_gen file to temp folder
		if err := ioutil.WriteFile(custom, src, 0644); err != nil {
			panic(err)
		}

		// add it to the command
		run = append(run, custom)
	}

	// write gen.go template to temp folder
	main := filepath.Join(tempDir, "main.go")

	err = ioutil.WriteFile(main, []byte(gentemplate), 0644)

	if err != nil {
		panic(err)
	}

	run = append(run, main)

	// run new gen

	var out bytes.Buffer
	var outerr bytes.Buffer

	cmd := exec.Command("go", run...)
	cmd.Stdout = &out
	cmd.Stderr = &outerr

	if err := cmd.Run(); err != nil {
		log.Println(outerr.String())
		panic(err)
	}

	if out.Len() > 0 {
		log.Println(out.String())
	}
}

const gentemplate = `package main

import (
	"github.com/clipperhouse/gen/typewriter"
	_ "github.com/clipperhouse/gen/typewriters/genwriter"
	_ "github.com/clipperhouse/gen/typewriters/container"
)

func main() {
	app, err := typewriter.NewApp("+gen")
	
	if err != nil {
		panic(err)
	}

	if err := app.WriteAll(); err != nil {
		panic(err)
	}
}
`
