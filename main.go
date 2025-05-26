package main

import (
	"fmt"
	"os"
	"sort"
)

// We let the program build a dynamic array of commands.  It becomes
// easy to add a new function then "register" it; the main() parser
// just works.

// We need the function to call and a help string
type fn struct {
	f    func()
	help string
}

// This is where we store all the commands
var fnmap = map[string]fn{}

// Register a new command
//
//	eg
//	   var _ bool = register_fn("help", help,"This help message")
func register_fn(cmd string, f func(), help string) bool {
	fnmap[cmd] = fn{f: f, help: help}
	return true
}

// Global variable for simplicity; the functions can read them
var Args []string

func main() {
	Args = os.Args

	// If no paramter is passed, default to "status"
	// otherwise split command line args
	var cmd string = "status"
	if len(Args) > 1 {
		cmd = Args[1]
		Args = Args[2:]
	}

	val, ok := fnmap[cmd]

	// If this is a defined command, call it
	if ok {
		val.f()
	} else {
		fmt.Println("Unknown command " + cmd)
		fmt.Println()
		help()
	}
}

func help() {
	fmt.Println("Command options:")

	// Is there an easier way to create a sorted key?

	keys := make([]string, 0, len(fnmap))
	for k := range fnmap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("  %10s: %s\n", k, fnmap[k].help)
	}
}

var _ bool = register_fn("help", help, "This help message")
