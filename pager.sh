#!/bin/bash
# Script wrapper for cleanpg

PROG="go run . -v -o output.html"

echo -e "Enter the URL: \c"
read URL

`$PROG $URL` && open ./output.html
