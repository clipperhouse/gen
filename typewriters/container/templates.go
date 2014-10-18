package container

import (
	"github.com/clipperhouse/gen/typewriter"
)

var templates = typewriter.TemplateSet{
	"List": list,
	"Ring": ring,
	"Set":  set,
}
