# gopop - Wrapper arround the Poppler CLI

This is a wrapper around the Poppler CLI tools. It is mostly a toy project to develop something
open source. Feel free to contribute or use it but don't expect high availability support or maintenance
(at least for the time being).

## Why wrap the CLI?

There are [a couple of modules integrating Poppler via C bindings](https://pkg.go.dev/search?q=poppler&m=).
This project aims to be an alternative that makes use of parsing the CLI output of Poppler commands.
This may not be the most advanced approach but it has the advantage of being cross-platform and
generally easier to implement (at least for me).
