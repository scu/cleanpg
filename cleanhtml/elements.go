// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cleanhtml

import "strings"

// nodeElements holds html.ElementNode objects
type nodeElements struct {
	attributes []string
	style      string
}

var encounteredBodyElement bool = false
var encounteredFirstH1Element bool = false

// isElementRenderable determines if the key "node"
// exists in the renderableHTML map
func isElementRenderable(node string) bool {
	lcaseTag := strings.ToLower(node)
	var doRender bool = false

	// Is it in the map
	if _, ok := renderableHTML[lcaseTag]; ok {
		doRender = true
	}

	// Skip link rendering
	if lcaseTag == "a" && doRender && !renderLinks {
		return false
	}

	// Special processing directives for "canonical mode"
	// which indicates only body & div elements are to be
	// rendered until the first h1 tag is encountered
	if renderCanonicalMode && doRender {

		if lcaseTag == "body" {
			encounteredBodyElement = true
		}

		if encounteredBodyElement &&
			!encounteredFirstH1Element &&
			lcaseTag != "body" {
			doRender = false
		}

		if lcaseTag == "h1" {
			encounteredFirstH1Element = true
			doRender = true
		}
	}

	return doRender
}

// isElementAttributeRenderable determines if they key "attr"
// exists in renderableHTML["node"]
func isElementAttributeRenderable(node string, attr string) bool {
	// Assume there are uppercase tags out there <Html>, <HTML> etc
	lcNode := strings.ToLower(node)

	// Check if the node exists
	if _, ok := renderableHTML[lcNode]; !ok {
		return false
	}

	// Loop through attributes
	if renderableHTML[lcNode].attributes != nil {
		for _, v := range renderableHTML[lcNode].attributes {
			if v == attr {
				return true
			}
		}
	}

	return false
}

// Void elements (can't have any contents)
// See section 12.1.2 of HTML reference
var voidElements = map[string]bool{
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"keygen": true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
}

// Elements and associated attributes to render.
// Specification at https://developer.mozilla.org/en-US/docs/Web/HTML/Element
var renderableHTML = map[string]nodeElements{
	// Main root
	"html": {
		style: `
		margin: auto;
		height: 100%;
		display: table;
		background: #d6dede;
		`,
	},

	// Document metadata
	"head":  {},
	"title": {},

	// Sectioning root
	"body": {
		style: `
		margin: 0 auto;
		padding-left: 20px;
		padding-right: 20px;
		height: 100%;
		font: 115% 'PT Sans', 'Helvetica', sans-serif;
		max-width: 800px;
		color: #555753; 
		background: #fff; 
		display: table-cell;
		vertical-align: middle;
		`,
	},
	"div":  {},
	"span": {},

	// Content sectioning
	"h1": {
		style: `
		font-size: 175%;
		margin-top: 40px;
		`,
	},
	"h2": {
		style: `
		font-size: 145%;
		margin-top: 30px;
		`,
	},
	"h3": {
		style: `
		font-size: 130%;
		margin-top: 20px;
		`,
	},
	"h4": {},
	"h5": {},
	"h6": {},

	// Text content
	"p":          {},
	"blockquote": {},
	"pre": {
		style: `font-family: Menlo, monospace;
		font-size: 0.875rem;`,
	},
	"code": {
		style: `font-family: Menlo, monospace;
		word-spacing: -0.3em;
		font-size: 0.875rem;`,
	},

	// Inline text semantics
	"a": {
		attributes: []string{
			"href",
		},
	},
	"b":  {},
	"em": {},
	"i":  {},
	"br": {},

	// Image and multimedia
	// TODO: option to catalog images
	//"img": {"src"},

	// Table content
	"caption":  {},
	"col":      {},
	"colgroup": {},
	"table":    {},
	"tbody":    {},
	"td":       {},
	"tfoot":    {},
	"th":       {},
	"thead":    {},
	"tr":       {},
}
