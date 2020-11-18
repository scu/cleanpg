// Copyright 2020 Scott Underwood.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Package cleanhtml provides read -> parse -> filter -> render
// capability to the cleanpg utility.
package cleanhtml

import (
	"io/ioutil"
	"net/http"

	"github.com/scu/cleanpg/logger"
)

// ReadHTML reads a web page and returns a string
// containing the unfiltered contents and an error
func ReadHTML(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		logger.Write(logger.FATAL, "Could not get url [%s]: %s", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	// read html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Write(logger.FATAL, "Could not read bytes from [%s]: %s", url, err)
		return nil, err
	}

	return html, nil
}
