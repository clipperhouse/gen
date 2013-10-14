package main

func getTemplates() (templates map[string]string) {
	templates = make(map[string]string)

	templates["header"] = `// {{.Command}}
// this file was auto-generated using github.com/clipperhouse/gen
// {{.Generated}}

`
	templates["package"] = `package {{.Package}}

`
	templates["type"] = `type {{.Plural}} []*{{.Singular}}

`
	templates["all"] = `func ({{.PluralLocal}} {{.Plural}}) All(fn func({{.SingularLocal}} *{{.Singular}}) bool) bool {
	if fn == nil {
		return true
	}
	for _, m := range {{.PluralLocal}} {
		if !fn(m) {
			return false
		}
	}
	return true
}
`
	templates["any"] = `func ({{.PluralLocal}} {{.Plural}}) Any(fn func({{.SingularLocal}} *{{.Singular}}) bool) bool {
	if fn == nil {
		return true
	}
	for _, m := range {{.PluralLocal}} {
		if fn(m) {
			return true
		}
	}
	return false
}
`
	templates["count"] = `func ({{.PluralLocal}} {{.Plural}}) Count(fn func({{.SingularLocal}} *{{.Singular}}) bool) (result int) {
	if fn == nil {
		return len({{.PluralLocal}})
	}
	for _, m := range {{.PluralLocal}} {
		if fn(m) {
			result++
		}
	}
	return result
}
`
	templates["where"] = `func ({{.PluralLocal}} {{.Plural}}) Where(fn func({{.SingularLocal}} *{{.Singular}}) bool) (result {{.Plural}}) {
	for _, m := range {{.PluralLocal}} {
		if fn == nil || fn(m) {
			result = append(result, m)
		}
	}
	return result
}
`
	return
}
