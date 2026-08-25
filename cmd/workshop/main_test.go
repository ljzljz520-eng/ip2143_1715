package main

import "testing"

func TestCLIExample(t *testing.T) {
	var entry func(string, string) error = Run
	_ = entry
}
