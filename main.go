package main

import (
	"bitbucket.org/pkg/inflect"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

type Values struct {
	Package   string
	Singular  string
	Plural    string
	Receiver  string
	Loop      string
	Pointer   string
	Generated string
	Command   string
	FileName  string
}

func main() {
	t := getTemplates()
	v := getValues()
	writeTemplates(t, v)
}

var arg = regexp.MustCompile(`([a-zA-Z]+)\.(\*?)([a-zA-Z]+)`)

func getValues() (v *Values) {
	matches := arg.FindStringSubmatch(os.Args[1])

	if matches == nil {
		log.Fatalln("The first argument must be in the form of package.TypeName")
	}

	pkg := matches[1]
	ptr := matches[2]
	typ := inflect.Singularize(matches[3])

	return &Values{
		Package:   pkg,
		Singular:  typ,
		Plural:    inflect.Pluralize(typ),
		Receiver:  string(strings.ToLower(typ)[0]),
		Loop:      "_" + string(strings.ToLower(typ)[0]),
		Pointer:   ptr,
		Generated: time.Now().UTC().Format(time.RFC1123),
		Command:   strings.Join(os.Args, " "),
		FileName:  strings.ToLower(typ) + "_gen.go",
	}
}

func writeTemplates(templates map[string]string, v *Values) {
	f, err := os.Create(v.FileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for key, val := range templates {
		t := template.Must(template.New(key).Parse(val))
		t.Execute(f, v)
	}
}
