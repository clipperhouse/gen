package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func main() {
	// read gen custom imports file
	custom, err := ioutil.ReadFile("_gen.go")
	if err != nil {
		// maybe silently fail here? some people may not use this feature at all
		// log.Println("no custom _gen.go for imports found")
	}
	// create ._gen dir
	createdDir := true
	err = os.Mkdir("._gen", 0777)
	if err != nil {
		createdDir = false
		log.Println(err)
	}
	// panic safe dir removal?
	defer func() {
		// don't blow away a dir we didn't create
		if createdDir {
			err = os.RemoveAll("._gen")
			if err != nil {
				panic(err)
			}
		}
	}()
	// minimal compiling file if none provided
	if len(custom) == 0 {
		custom = []byte("package main")
	}
	// write custom_gen file to ._gen folder
	err = ioutil.WriteFile("._gen/gen_custom.go", custom, 0644)
	if err != nil {
		panic(err)
	}
	// write gen.go template to ._gen folder
	err = ioutil.WriteFile("._gen/gen.go", []byte(gentemplate), 0644)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	var outerr bytes.Buffer
	// build new gen
	cmd := exec.Command("go", "build", "-o", "._gen/gen", "._gen/gen.go", "._gen/gen_custom.go")
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
	// run new gen
	cmd = exec.Command("._gen/gen", os.Args...)
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
	"github.com/clipperhouse/typewriter"
	_ "github.com/clipperhouse/typewriters/genwriter"
	_ "github.com/clipperhouse/typewriters/container"
)

func main() {
	app, err := typewriter.NewApp("+gen")
	if err != nil {
		panic(err)
	}

	app.WriteAll()
}
`
