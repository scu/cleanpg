// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package main provides an entry point for the cleanpg utility
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/scu/cleanpg/cleanhtml"
	"github.com/scu/cleanpg/logger"
	"github.com/scu/flagplus"
)

func usage() {
	fmt.Println(fs.Usage())
}

func init() {
	fs = flagplus.NewFlagSet("cleanpg")
	fs.FlagSetDescription("Utility for rendering text-readable versions of HTML pages.")

	// Add flags
	fs.AddFlag("verbose", "v", "Print extra debugging information to stderr")
	fs.AddFlag("help", "h", "Help")
	fs.AddFlag("nocanon", "c", "Do not attempt to render canonically")
	fs.AddFlag("nostyle", "n", "Do not render embedded style")
	fs.AddFlag("nolinks", "l", "Do not render links")
	fs.AddStringFlag("output", "o", "Write output to `file.html`", "out.html")
	fs.AddStringFlag("save", "s", "Save source document as `file.html`", "")

	fs.Parse()
}

var fs *flagplus.FlagSet

func main() {
	// Call cleanpgMain in a separate function
	// so that it deferred statements run before exit
	exitCode := cleanpgMain()
	os.Exit(exitCode)
}

func cleanpgMain() int {

	var outFile *os.File

	// Set up logging
	logger.Truncate()

	// FLAG "help"
	help, err := fs.Get("help")
	if err != nil {
		panic(err)
	}
	if help {
		usage()
		return 0
	}

	// FLAG "verbose"
	printVerbose, err := fs.Get("verbose")
	if err != nil {
		panic(err)
	}
	if printVerbose {
		logger.LogToStderr(true)
	}

	// FLAG "output"
	outputFile, err := fs.GetString("output")
	if err != nil {
		panic(err)
	}
	if outputFile != "" {
		// Verify is .html extension
		if filepath.Ext(outputFile) != ".html" {
			logger.Write(logger.FATAL, "file [%s] must have .html extension", outputFile)
			return 1
		}

		var err error
		// Create & open the file
		outFile, err = os.Create(outputFile)
		if err != nil {
			logger.Write(logger.FATAL, "could not open [%s]: %s", outputFile, err)
			return 1
		}
	}

	// Get url from argument
	args := fs.GetArgs()
	urlToClean := args[0]
	if urlToClean == "" {
		fmt.Fprintf(os.Stderr, "Missing URL\n")
		usage()
		return 1
	}

	logger.Write(logger.INFO, "reading data from URL=%s", urlToClean)

	sourceData, err := cleanhtml.ReadHTML(urlToClean)
	if err != nil {
		logger.Write(logger.FATAL, "Cannot read [%s]: %s", urlToClean, err)
		return 1
	}

	// FLAG "save"
	saveFile, err := fs.GetString("save")
	if err != nil {
		panic(err)
	}
	// Save a copy in source.html
	if saveFile != "" {
		// Verify is .html extension
		if filepath.Ext(saveFile) != ".html" {
			logger.Write(logger.FATAL, "file [%s] must have .html extension", saveFile)
			return 1
		}
		svFile, err := os.Create(saveFile)
		if err != nil {
			logger.Write(logger.FATAL, "could not open save file [%s]: %s", saveFile, err)
			return 1
		}
		defer svFile.Close()
		logger.Write(logger.INFO, "saving a copy of the source document to %s", saveFile)
		fmt.Fprintf(svFile, "%s", sourceData)
	}

	// FLAG "nocanon"
	nocanon, err := fs.Get("nocanon")
	if err != nil {
		panic(err)
	}
	if !nocanon {
		// Canonical is default
		cleanhtml.SetPostH1Render(true)
		logger.Write(logger.INFO, "processing body elements after first <h1> tag")
	}

	// FLAG "nostyle"
	noStyle, err := fs.Get("nostyle")
	if err != nil {
		panic(err)
	}
	if noStyle {
		cleanhtml.SetStyleRender(false)
		logger.Write(logger.INFO, "skipping automatic tag-level style embedding")
	}

	// FLAG "nolinks"
	noLinks, err := fs.Get("nolinks")
	if err != nil {
		panic(err)
	}
	if noLinks {
		cleanhtml.SetLinksRender(false)
		logger.Write(logger.INFO, "not rendering links")
	}

	// Create the cleanly-formatted page
	cleanData, err := cleanhtml.CleanHTML(sourceData)
	if err != nil {
		logger.Write(logger.FATAL, "Could not clean [%s]: %s", urlToClean, err)
		return 1
	}

	// Write to designated output
	fmt.Fprintf(outFile, "%s", cleanData)
	fmt.Printf("Document rendered to %q\n", outputFile)
	logger.Write(logger.INFO, "Document from %q rendered to %q", urlToClean, outputFile)

	return 0

}
