##What’s this?

Typewriter is a new package to enable pluggable, type-based codegen for Go. The envisioned use case is for generics-like functionality. This package forms the underpinning of (the next version of) [gen](https://github.com/clipperhouse/gen/tree/master/typewriter).

Usage is analogous to how codecs work with Go’s [image](http://golang.org/pkg/image/) package, or database drivers in the [sql](http://golang.org/pkg/database/sql/) package.

    import (
        // main package
    	“github.com/clipperhouse/gen/typewriter”
    	
    	// any number of typewriters 
    	_ “github.com/clipperhouse/gen/typewriters/container”
    	_ “github.com/clipperhouse/gen/typewriters/genwriter”
    )
    
    func main() {
    	app, err := typewriter.NewApp(”+gen”)
    	if err != nil {
    		panic(err)
    	}
    
    	app.WriteAll()
    }

Individual [typewriters](https://github.com/clipperhouse/gen/tree/master/typewriters) register themselves to the “parent” package via their init() functions.

This is new and in-progress. Feedback is welcome.
