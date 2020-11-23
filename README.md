# Utility cleanpg

cleanpg is a tool for rendering a source HTML document into a more human-readable format.

By default, the document is written to `out.html` in the current directory. To override, use the `-o file` (or `--output file`) command line flag. Note: file extension must be .html.

| Original | Rendered |
| ------------------ | ------------------ |
|![Before](htmlb4.png)| ![After](htmlafter.png) |

The utility defaults to **canonical mode** which applies specific
assumptions to improve readability, such as skipping over elements between the `<body>` tag and the first `<h1>` tag. Canonical mode may be turned
off by using the `-c` (or `--nocanon`) command line flag.

Tag-level styles are embedded for readability. For example, `<h1 style="font-size: 175%;margin-top: 40px;">` is embedded automatically for each H1 element. **Disable this default behavior** by using the `-n` (or `--nostyle`) command line flag.

Links are rendered by default. To skip links, use the `-l` (or `--nolinks`) command line flag.

### Disclaimer:
cleanpg re-renders document ("page") layouts and content for experimental use only. Use of these altered pages may not be used for re-publishing, circumventing content protection schemes, or in any manner which violates copyright law.

## Getting Started

### Installation:
```
go get github.com/scu/cleanpg
```

### Basic usage:
```
cleanpg url
```

### Example:
```
cleanpg http://example.org
```

## Command-line options
```
Utility for rendering text-readable versions of HTML pages.
Usage:
  cleanpg [-h|c|l|n|o file.html|s file.html|v]
Options:
  -h, --help 
     Help
  -c, --nocanon 
     Do not attempt to render canonically
  -l, --nolinks 
     Do not render links
  -n, --nostyle 
     Do not render embedded style
  -o, --output file.html
     Write output to file.html (default=out.html)
  -s, --save file.html
     Save source document as file.html
  -v, --verbose 
     Print extra debugging information to stderr
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Versioning

[SemVer](http://semver.org/) is used for versioning. For the versions available, see the [tags on this repository](https://github.com/scu/cleanpg/tags). 

## Authors

* **Scott Underwood** - *Initial work* - [Scott Underwood](https://github.com/scu)

## License
[MIT](https://choosealicense.com/licenses/mit/)