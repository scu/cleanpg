// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

/*
Package cleanhtml provides a toolset for reading source HTML documents
and attempting to render them into more human-readable output.

Although this package is meant to be consumed by the cleanpg utility
(http://github.com/scu/cleanpg) it may be useful in other applications.

	url := "http://example.com"
	sourceData, err := cleanhtml.ReadHTML(url)
	if err != nil {
		errStr := fmt.Sprintf("Could not read document at %q: %s", url, err)
		panic(errStr)
	}

	cleanData, err := cleanhtml.CleanHTML(sourceData)
	if err != nil {
		errStr := fmt.Sprintf("Could not transform data: %s", err)
		panic(errStr)
	}

Disclaimer: this library outputs a document layout and content different
than the original page designer. Use of these re-rendered documents are
not intended for re-publishing, circumventing content protection mechanisms
or violate the copyright of the original content authors. */
package cleanhtml
