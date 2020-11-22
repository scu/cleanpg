// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cleanhtml

import (
	"bytes"
	"io"

	"github.com/scu/cleanpg/logger"
	"golang.org/x/net/html"
)

var renderCanonicalMode bool = false

// SetPostH1Render sets flag indicating whether
// the renderer will process BODY elements until the
// first H1 tag is reached
func SetPostH1Render(flag bool) {
	renderCanonicalMode = flag
}

var renderStyle bool = true

// SetStyleRender sets flag indicating whether
// the renderer embeds tag-level styles automatically
// [default = true]
func SetStyleRender(flag bool) {
	renderStyle = flag
}

var renderLinks bool = true

// SetLinksRender sets flag indicating whether
// links <a... href...> will be rendered
// [default = true]
func SetLinksRender(flag bool) {
	renderLinks = flag
}

// CleanHTML provides a rendered HTML document.
// It accepts document data (normally through cleanhtml.ReadHTML),
// parses and renders the data through a set of filters to produce
// readable HTML output, which is returned as a string.
func CleanHTML(data []byte) (string, error) {
	// Parse the document
	docNodes, err := html.Parse(bytes.NewReader(data))
	if err != nil {
		logger.Write(logger.FATAL, "Could not parse HTML: %s", err)
		return "", err
	}

	var buf bytes.Buffer
	w := io.Writer(&buf)
	render(w.(writer), docNodes)
	return buf.String(), nil
}
