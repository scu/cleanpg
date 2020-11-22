// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package main provides an entry point for the cleanpg utility
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/scu/cleanpg/cleanhtml"
	"github.com/scu/cleanpg/logger"
)

// Command line flags usage map
// So usage descriptions can be shared across long and short options
var usageMap = map[string]string{
	"help":    "Help",
	"verbose": "Print extra debugging information",
	"output":  "HTML file to render to [default=stdout]",
	"save":    "Save a copy of the source HTML document",
	"posth1":  "Render body elements after first h1 tag",
	"nostyle": "Do not automatically render tag-level embedded styles",
	"nolinks": "Do not render links",
}

// Command line flags
var (
	helpPtr    = flag.Bool("help", false, usageMap["help"])
	verbosePtr = flag.Bool("verbose", false, usageMap["verbose"])
	outputPtr  = flag.String("output", "", usageMap["output"])
	savePtr    = flag.String("save", "", usageMap["save"])
	posth1Ptr  = flag.Bool("posth1", false, usageMap["posth1"])
	noStylePtr = flag.Bool("nostyle", false, usageMap["nostyle"])
	noLinksPtr = flag.Bool("nolinks", false, usageMap["nolinks"])
)

// Holds the fd for the output [default=stdout]
var outFile *os.File = os.Stdout

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: cleanpg [flags] url\n")
	flag.PrintDefaults()
}

func init() {
	// Add optional short arguments to flag map
	flag.BoolVar(helpPtr, "?", false, usageMap["help"])
	flag.BoolVar(verbosePtr, "v", false, usageMap["verbose"])
	flag.StringVar(outputPtr, "o", "", usageMap["output"])
	flag.StringVar(savePtr, "s", "", usageMap["save"])
	flag.BoolVar(posth1Ptr, "p", false, usageMap["posth1"])
	flag.BoolVar(noStylePtr, "n", false, usageMap["nostyle"])
	flag.BoolVar(noLinksPtr, "l", false, usageMap["nolinks"])
}

func main() {
	// Call cleanpgMain in a separate function
	// so that it deferred statements run before exit
	exitCode := cleanpgMain()
	os.Exit(exitCode)
}

func cleanpgMain() int {

	// Set up logging
	logger.Truncate()

	// Flags
	flag.Usage = usage
	flag.Parse()

	// Print help and exit
	if *helpPtr {
		usage()
		return 0
	}

	// Get url from argument
	urlToClean := flag.Arg(0)
	if urlToClean == "" {
		fmt.Fprintf(os.Stderr, "Missing URL\n")
		usage()
		return 1
	}

	// Optional flag: print extra data to stderr
	if *verbosePtr {
		logger.LogToStderr(true)
	}

	if *outputPtr != "" {
		// Verify is .html extension
		if filepath.Ext(*outputPtr) != ".html" {
			logger.Write(logger.FATAL, "file [%s] must have .html extension", *outputPtr)
			return 1
		}

		var err error
		// Create & open the file
		outFile, err = os.Create(*outputPtr)
		if err != nil {
			logger.Write(logger.FATAL, "could not open [%s]: %s", *outputPtr, err)
			return 1
		}
	}

	logger.Write(logger.INFO, "reading data from URL=%s", urlToClean)

	sourceData, err := cleanhtml.ReadHTML(urlToClean)
	if err != nil {
		logger.Write(logger.FATAL, "Cannot read [%s]: %s", urlToClean, err)
		return 1
	}

	// Save a copy in source.html
	if *savePtr != "" {
		// Verify is .html extension
		if filepath.Ext(*savePtr) != ".html" {
			logger.Write(logger.FATAL, "file [%s] must have .html extension", *savePtr)
			return 1
		}
		svFile, err := os.Create(*savePtr)
		if err != nil {
			logger.Write(logger.FATAL, "could not open save file [%s]: %s", *savePtr, err)
			return 1
		}
		defer svFile.Close()
		logger.Write(logger.INFO, "saving a copy of the source document to %s", *savePtr)
		fmt.Fprintf(svFile, "%s", sourceData)
	}

	// Optional flag: process body data only after h1
	// flag is encountered (exception is <div> elements since
	// <h1> may be encapsulated)
	if *posth1Ptr {
		cleanhtml.SetPostH1Render(true)
		logger.Write(logger.INFO, "processing body elements after first <h1> tag")
	}

	// Optional flag: do not embed tag-level style.
	if *noStylePtr {
		cleanhtml.SetStyleRender(false)
		logger.Write(logger.INFO, "skipping automatic tag-level style embedding")
	}

	// Optional flag: do not render links.
	if *noLinksPtr {
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
	logger.Write(logger.INFO, "created a clean version of %s", urlToClean)

	return 0

}
