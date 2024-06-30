package main

import (
	"bytes"
	"os"
	"testing"
)

const (
	inputFile = "./testdata/test.md"
	resultFile = "test.md.html"
	goldenFile = "./testdata/test.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)

	expected, err := os.ReadFile(goldenFile)

	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Error("Result content doesn't match golden file")
	}
}