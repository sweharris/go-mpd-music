package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func die(v ...any) {
	fmt.Fprintln(os.Stderr, "\nError:")
	fmt.Fprintf(os.Stderr, "  %v\n", v...)
	os.Exit(1)
}

// Convert a string number which represents seconds into a m:ss value
func secs(t_str string) string {
	if t_str == "" {
		t_str = "0"
	}
	t, e := time.ParseDuration(t_str + "s")
	if e != nil {
		die(e)
	}
	t = t.Round(time.Second)
	m := t / time.Minute
	t -= m * time.Minute
	s := t / time.Second
	return fmt.Sprintf("%d:%02d", m, s)
}

func truefalse(s string) bool {
	s = strings.ToLower(s)
	if s == "true" || s == "on" {
		return true
	} else if s == "flase" || s == "off" {
		return false
	}
	die("Not true/false/on/off")
	return false
}

func show_truefalse(s string) string {
	if s == "0" {
		return "false"
	} else if s == "1" {
		return "true"
	} else {
		return s
	}
}
