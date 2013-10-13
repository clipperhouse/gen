package main

import (
	"bitbucket.org/pkg/inflect"
	"os"
	"strings"
	"text/template"
	"time"
)

type Values struct {
	Package       string
	Singular      string
	SingularLocal string
	Plural        string
	PluralLocal   string
	Generated     string
	Command       string
	FileName      string
}

func main() {
	t := getTemplates()
	v := getValues()
	writeTemplates(t, v)
}

func getValues() (v *Values) {
	a := strings.Split(os.Args[1], ".")
	model := inflect.Singularize(a[1])
	models := inflect.Pluralize(model)
	return &Values{
		Package:       a[0],
		Singular:      model,
		SingularLocal: strings.ToLower(model),
		Plural:        models,
		PluralLocal:   strings.ToLower(models),
		Generated:     time.Now().UTC().Format(time.RFC1123),
		Command:       strings.Join(os.Args, " "),
		FileName:      strings.ToLower(model) + "_gen.go",
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
