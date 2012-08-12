title: It's November
template: blog-entry.html
---
Combining a source file with a template to produce an output file is simple
enough, but where does that output file get placed? To begin, all output files
go under the `-output-path` directory. Their relative location underneath that
path is determined by their `output` metadata key.

If no `output` metadata is specified, grender will use the complete relative
path of the source file, minus its extension.

