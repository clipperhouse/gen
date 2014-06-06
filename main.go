package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
)

func main() {
	// read gen custom imports file
	custom, err := ioutil.ReadFile("_gen.go")
	if err != nil {
		// maybe silently fail here? some people may not use this feature at all
		// log.Println("no custom _gen.go for imports found")
	}

	// minimal compiling file if none provided
	if len(custom) == 0 {
		custom = []byte("package main")
	}

	caller := path.Base(os.Args[0])
	tempDir, err := ioutil.TempDir("", caller)

	if err != nil {
		log.Println(err)
	}

	defer os.RemoveAll(tempDir)

	// write custom_gen file to temp folder
	err = ioutil.WriteFile(path.Join(tempDir, "gen_custom.go"), custom, 0644)
	if err != nil {
		panic(err)
	}

	// write gen.go template to temp folder
	err = ioutil.WriteFile(path.Join(tempDir, "gen.go"), []byte(gentemplate), 0644)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	var outerr bytes.Buffer

	// run new gen
	cmd := exec.Command("go", "run", path.Join(tempDir, "gen.go"), path.Join(tempDir, "gen_custom.go"))
	cmd.Stdout = &out
	cmd.Stderr = &outerr
	err = cmd.Run()
	if err != nil {
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

	app.WriteAll()
}
`
