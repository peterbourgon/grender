# Grender

Grender is a static site generator. It combines source files and template files
to produce a website.

## File types

### Source files

A source file contains the content of a page, plus some metadata. Content can be
raw HTML, or some renderable markup language, like Markdown. Metadata can be
explicitly specified, or can be deduced from other properties of the source
file, like its basename (the filename without the extension).

```
template: basic.template
title: My sample page
---
This is the **Markdown content** of my sample source file.
```

### Template files

A template file contains markup used to render an HTML file.

```
<html>
<head><title>{{title}}</title></head>
<body>
<h1>Welcome to my sample template file!</h1>
{{{content}}}
</body>
</html>
```


## Concepts

### Metadata

Metadata occurs at the beginning of a source file. It's delimited (terminated)
by the `-metadata-delimiter`. It's parsed as YAML. Grender reads (consumes) a
few specific pieces of metadata as part of its processing, but all others are
passed to the template in the context.

### Source file content

Source file content (everything after the `-metadata-delimiter`) is rendered as
Markdown, and placed in the context under the `-content-key`.

### Context (object)

The context is all of the information that's given to a template, so that the
template can be rendered to an HTML file. The context is primarily populated by
the metadata portion of a given source file. Rendered source file content is
provided under the `-content-key`. It should probably be referenced in the
template with three sets of curly braces (eg. `{{{content}}}`) so that HTML tags
aren't escaped.

### Special case: blog entries

Blog entries require additional, special treatment from the static site
generator. Blog entries are organized and accessed by date. It should be
possible to build an index page of blog entries, or some subset of them. At
least one page will want to render "the most recent" blog entry. When viewing a
blog entry it should be possible to navigate to the "next" and "previous" blog
entries.

To support these requirements, grender implements a concept called the
**index**.

### The index

To join the index, a source page should define a map called `index` in its
metadata. That map should contain at least one string called `key`, containing a
unique value. Call this `index` map an **index-tuple**.

As grender analyzes source files, it collects all of these index-tuples. The
tuples are first organized into groups according to their `type` (if not
present, `-default-index-type` is used). Then, each group sorts its tuples
according to the their `key` (decreasing). This aggregate, ordered data
structure is provided to every template rendering context as the `index`.

To support the concept of showing "the most recent" blog entries on a certain
page, the rendered content of the first `-index-content-count` index-tuples for
each `type` is provided in the appropriate index-tuples, under the
`-content-key`.

As a convenience, source files whose basenames match the pattern
`YYYY-MM-DD-title-text` are interpreted as blog entries, and automatically get
an `index` metadata key, populated as follows:

```
index:
	type: [-default-index-type]
	key: YYYY-MM-DD-title-text
	year: YYYY
	month: MM
	day: DD
	title: title text
	url: [-blog-path]/YYYY-MM-DD-title-text.[-output-extension]
```

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

* `[-source-path]/foo/bar.md` will generate
  `[-output-path]/foo/bar.[-output-extension]`

As a special case, source files whose basenames match the pattern
`YYYY-MM-DD-title-text` are interpreted as blog entries, and automatically get
an `output` metadata key, populated with
`[-output-path]/[-blog-path]/YYYY-MM-DD-title-text.[-output-extension]`.
