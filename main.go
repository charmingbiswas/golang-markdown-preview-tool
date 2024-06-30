package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Markdown Preview Tool</title>
		</head>
		<body>
	`

	footer = `
		</body>
		</html>
	`
)

func main() {
	// parse the flags
	filename := flag.String("file", "", "Name of the markdown file")
	flag.Parse()

	// if user did not input name of the file, show usage

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


func run(filename string) error {
	// read all the data from input file and check for errors

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(fileContent)

	outName := fmt.Sprintf("%s.html", filepath.Base(filename))

	fmt.Println(outName)
	openBrowser(outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	
	// this block of code generates a valid block of HTML
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Create buffer to write to file
	var buffer bytes.Buffer
	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outname string, htmldata []byte) error {
	// this function will store the contents of the buffer of bytes above as HTML file
	// 0644 refers to read and write permissions for a file by the OS
	return os.WriteFile(outname, htmldata, 0644)
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}