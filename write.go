package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

// set up temp directory under current directory
// make sure to defer os.RemoveAll() in caller
func getTempDir() (string, error) {
	caller := filepath.Base(os.Args[0])
	wd, _ := os.Getwd()
	return ioutil.TempDir(wd, caller)
}

var tmpl = template.Must(template.New("package").Parse(`package {{.Name}}
{{if gt (len .Imports) 0}}
import ({{range .Imports}}
	{{.}}{{end}}
)
{{end}}`))

const runBody string = `
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		log.Fatal(err)
	}

	if err := app.WriteAll(); err != nil {
		log.Fatal(err)
	}
}
`

const listBody string = `
func main() {
	app, err := typewriter.NewApp("+gen")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Installed typewriters (custom):")
	for _, tw := range app.TypeWriters {
		fmt.Println("  " + tw.Name())
	}
}
`
