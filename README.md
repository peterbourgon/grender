# Grender

Grender is a static site generator. It combines source files and template files
to produce a website.

[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/peterbourgon/grender.png
[2]: http://www.travis-ci.org/peterbourgon/grender

## File types

### Source files

A source file contains the content of a page, plus some [metadata](#metadata).
Content can be raw HTML, or some renderable markup language, like
[Markdown][markdown]. Metadata can be explicitly specified, or can be deduced
from other properties of the source file, like its basename (the filename
without the extension).

[markdown]: http://daringfireball.net/projects/markdown/syntax

```
template: basic.template
title: My sample page
---
This is the **Markdown content** of my sample source file.
```

### Template files

A template file contains markup and [Mustache][mustache] template tags, used to
render an HTML file.

[mustache]: http://github.com/hoisie/mustache

```
<html>
<head><title>{{title}}</title></head>
<body>
<h1>Welcome to my sample template file!</h1>
{{{content}}}
</body>
</html>
```

### Static files

A static file is copied verbatim to the `-output-path`.


## Concepts

### Metadata

Metadata occurs at the beginning of a source file. It's delimited (terminated)
by the `-metadata-delimiter`. It's parsed as YAML. Grender reads (consumes) a
few specific pieces of metadata as part of its processing, but all others are
passed to the template in the [context](#context-object).


### Source file content

Source file content (everything after the `-metadata-delimiter`) is rendered
according to the extension of the source file and placed in the
[context](#context-object) under the `-content-key`. An extension of `.md`
implies Markdown; any other extension implies raw data, ie. no rendering will
be performed.


### Context (object)

The context is all of the information that's given to a template, so that the
template can be rendered to an HTML file. The context is primarily populated by
the metadata portion of a given source file. Rendered source file content is
provided under the `-content-key`. It should probably be referenced in the
template with three sets of curly braces (eg. `{{{content}}}`) so that HTML tags
aren't escaped.


### Special case: blog entries

Blog entries are a special type of source file. They need to be available, in
whole or part, to templates that want to build lists of entries, or display the
most-recent-N entries entirely, with links to the previous and next page(s) of
entries.

To support these requirements, grender implements a concept called the
**index**.


### The index

A source page automatically becomes a member of the **index** if its name
matches a specific pattern: `YYYY-MM-DD-title-of-entry`. Based on this name, a
few metadata fields are automatically populated (but may be overridden). These
are best illustrated by example. A file named `2012-01-02-my-entry-name.md` will
have the following metadata populated:

```
title: My entry name
year: 2012
month: 01
day: 02
output: [-output-path]/[-blog-path]/YYYY-MM-DD-my-entry-name.[-output-extension]
url: [-blog-path]/2012-01-02-my-entry-name.[-output-extension]
sortkey: 2012-01-02-my-entry-name
```

During runtime, every source file that qualifies for the index will have its
(rendered) context object collected into a list, sorted by the sortkey. That
complete sorted list is provided to every source file specifying `index: true`
in its metadata, under the same `index` key (ie. replacing `true`).


### Output location

Combining a source file with a template to produce an output file is simple
enough, but where does that output file get placed? To begin, all output files
go under the `-output-path` directory. Their relative location underneath that
path is determined by their `output` metadata key. For example,

* `output: foo` will generate `[-output-path]/foo.[-output-extension]`

* `output: a/b/c/bar.baz` will generate
  `[-output-path]/a/b/c/bar.baz.[-output-extension]`

If no `output` metadata is specified, grender will use the complete relative
path of the source file, minus its extension. For example,

* `[-source-path]/foo/bar.md` implies `output: foo/bar` and and will generate
  `[-output-path]/foo/bar.[-output-extension]`
