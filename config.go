package main

import (
	"io"
	"os"
	"sync"
)

// lock for changing configs below
var mu = &sync.Mutex{}

// global state for output; useful with testing
var (
	defaultOut io.Writer = os.Stdout
	out        io.Writer = defaultOut
)

// with inspiration from http://golang.org/src/pkg/log/log.go?s=7258:7285#L229
func setOutput(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	out = w
}

func revertOutput() {
	setOutput(defaultOut)
}

// global state for custom imports file name; useful with testing
const defaultCustomName string = "_gen.go"

var customName string = defaultCustomName

func setCustomName(s string) {
	mu.Lock()
	defer mu.Unlock()
	customName = s
}

func revertCustomName() {
	setCustomName(defaultCustomName)
}
