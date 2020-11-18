// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package cleanhtml provides read -> parse -> filter -> render
// capability to the cleanpg utility.
package cleanhtml

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

func escape(w writer, s string) error {
	const escapedChars = "&'<>\"\r"

	i := strings.IndexAny(s, escapedChars)
	for i != -1 {
		if _, err := w.WriteString(s[:i]); err != nil {
			return err
		}
		var esc string
		switch s[i] {
		case '&':
			esc = "&amp;"
		case '\'':
			// "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
			esc = "&#39;"
		case '<':
			esc = "&lt;"
		case '>':
			esc = "&gt;"
		case '"':
			// "&#34;" is shorter than "&quot;".
			esc = "&#34;"
		case '\r':
			esc = "&#13;"
		default:
			panic("unrecognized escape character")
		}
		s = s[i+1:]
		if _, err := w.WriteString(esc); err != nil {
			return err
		}
		i = strings.IndexAny(s, escapedChars)
	}
	_, err := w.WriteString(s)
	return err
}

// isTextWhitespace returns true if text block runes
// are entirely whitespace
func isTextWhitespace(text string) bool {
	for _, v := range text {
		if !unicode.IsSpace(v) {
			return false
		}
	}

	return true
}

// cleanStyle accepts a string in "text" and removes
// newlines, carriage-returns, and tab characters
// so it can be inserted inline into the element tag
// style attribute
func cleanStyle(text string) string {
	var b strings.Builder
	b.Grow(len(text))

	for _, v := range text {
		if v != '\r' && v != '\n' && v != '\t' {
			b.WriteRune(v)
		}
	}

	return b.String()
}

// renderStartTag renders the start tag "\n<tag attr...>"
func renderStartTag(w writer, n *html.Node) error {
	// Begin element with a NL (for readability)
	if err := w.WriteByte('\n'); err != nil {
		return err
	}
	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil {
		return err
	}
	if _, err := w.WriteString(n.Data); err != nil {
		return err
	}

	// Add style attribute if present
	if renderStyle && renderableHTML[n.Data].style != "" {
		styleAttrib := fmt.Sprintf(" style=\"%s\"", renderableHTML[n.Data].style)
		if _, err := w.WriteString(cleanStyle(styleAttrib)); err != nil {
			return err
		}
	}

	// Render any attributes
	if err := renderAttributes(w, n); err != nil {
		return err
	}

	if voidElements[n.Data] {
		if n.FirstChild != nil {
			return fmt.Errorf("html: void element <%s> has child nodes", n.Data)
		}
		_, err := w.WriteString("/>")
		return err
	}
	if err := w.WriteByte('>'); err != nil {
		return err
	}

	return nil
}

// renderCloseTag renders the closing tag "</tag>"
func renderCloseTag(w writer, n *html.Node) error {
	closeTag := fmt.Sprintf("</%s>", n.Data)
	if _, err := w.WriteString(closeTag); err != nil {
		return err
	}
	return nil
}

// renderAttributes renders an html.ElementNode's attributes
func renderAttributes(w writer, n *html.Node) error {
	// Check attributes on html.ElementNode
	for _, a := range n.Attr {
		if isElementAttributeRenderable(n.Data, a.Key) {
			if err := w.WriteByte(' '); err != nil {
				return err
			}
			if a.Namespace != "" {
				if _, err := w.WriteString(a.Namespace); err != nil {
					return err
				}
				if err := w.WriteByte(':'); err != nil {
					return err
				}
			}
			keyVal := fmt.Sprintf("%s=\"%s\"", a.Key, a.Val)
			if _, err := w.WriteString(keyVal); err != nil {
				return err
			}
		}
	}
	return nil
}

// render is the main entry point for the rendering engine
func render(w writer, n *html.Node) error {
	// Render all nodes except ElementNode
	switch n.Type {
	case html.ErrorNode:
		return fmt.Errorf("cleanhtml: error node [%s]", n.Data)
	case html.TextNode:
		if !isTextWhitespace(n.Data) {
			escape(w, n.Data)
		}
		return nil
	case html.DocumentNode:
		// Starts here, render each node in doc nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := render(w, c); err != nil {
				return err
			}
		}
		return nil
	case html.ElementNode:
		// No operation, skip to below
	case html.CommentNode:
		// Do not render comments
		return nil
	case html.DoctypeNode:
		// Use our own doctype instead of source
		if _, err := w.WriteString("<!DOCTYPE html>"); err != nil {
			return err
		}
		return nil
	case html.RawNode:
		_, err := w.WriteString(n.Data)
		return err
	default:
		return errors.New("html: unknown node type")
	}

	// Determine if renderable
	renderElement := isElementRenderable(n.Data)

	if renderElement {
		if err := renderStartTag(w, n); err != nil {
			return err
		}
	}

	// Render child nodes.
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Don't render a TextNode if the parent element is unrenderable (i.e. <script>...</script>)
		if c.Type == html.TextNode && !isElementRenderable(c.Parent.Data) {
			continue
		}
		if err := render(w, c); err != nil {
			return err
		}
	}

	if renderElement {
		if err := renderCloseTag(w, n); err != nil {
			return err
		}
	}

	return nil
}

// writeQuoted writes s to w surrounded by quotes
func writeQuoted(w writer, s string) error {
	var q byte = '"'

	// Look for double quote, replace with escaped single
	if strings.Contains(s, `"`) {
		q = '\''
	}

	// Open quote
	if err := w.WriteByte(q); err != nil {
		return err
	}

	// Write string
	if _, err := w.WriteString(s); err != nil {
		return err
	}

	// Close quote
	if err := w.WriteByte(q); err != nil {
		return err
	}
	return nil
}
