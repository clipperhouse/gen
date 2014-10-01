package slice

import "github.com/clipperhouse/gen/typewriter"

var slice = &typewriter.Template{
	Text: `// {{.SliceName}} is a slice of type {{.Type}}. Use it where you would use []{{.Type}}.
type {{.SliceName}} []{{.Type}}
`,
}
