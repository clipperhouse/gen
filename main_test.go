package main

import (
	"strings"
	"testing"
)

type parseTest struct {
	text  string
	cmd   string
	force bool
	tail  int  //length
	err   bool //exists
}

func TestParseArgs(t *testing.T) {
	tests := []parseTest{
		parseTest{"gen", "", false, 0, false},
		parseTest{"gen -f", "", true, 0, false},
		parseTest{"gen yadda", "", false, 0, true}, // unknown command
		parseTest{"gen -bar", "", false, 0, true},  // tail is not ok
		parseTest{"gen add", "add", false, 0, false},
		parseTest{"gen add foo bar", "add", false, 2, false}, // tail is ok
		parseTest{"gen add -f", "add", true, 0, true},        // force is not ok
		parseTest{"gen get", "get", false, 0, false},
		parseTest{"gen get foo bar", "get", false, 2, false}, // tail is ok
		parseTest{"gen get -f", "get", true, 0, true},        // force is not ok
		parseTest{"gen help", "help", false, 0, false},
		parseTest{"gen help foo bar", "help", false, 0, true}, // tail is not ok
		parseTest{"gen help -f", "help", true, 0, true},       // force is not ok
		parseTest{"gen list", "list", false, 0, false},
		parseTest{"gen list foo bar", "list", false, 0, true}, // tail is not ok
		parseTest{"gen list -f", "list", true, 0, true},       // force is not ok
		parseTest{"gen watch", "watch", false, 0, false},
		parseTest{"gen watch foo bar", "watch", false, 0, true}, // tail is not ok
		parseTest{"gen watch -f", "watch", true, 0, false},      // force is ok
	}

	for i, test := range tests {
		cmd, force, tail, err := parseArgs(strings.Split(test.text, " "))
		if cmd != test.cmd {
			t.Errorf("tests[%d]: cmd should be %q, got %q", i, test.cmd, cmd)
		}
		if force != test.force {
			t.Errorf("tests[%d]: force should be %v, got %v", i, test.force, force)
		}
		if len(tail) != test.tail {
			t.Errorf("tests[%d]: len(tail) should be %v, got %v", i, test.tail, len(tail))
		}
		if (err != nil) != test.err {
			t.Errorf("tests[%d]: err existence should be %v, got %v", i, test.err, err)
		}
	}
}
