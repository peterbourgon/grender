# Grender

Grender is a static site generator. It combines source files with metadata to
produce a website.

[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/peterbourgon/grender.png
[2]: http://www.travis-ci.org/peterbourgon/grender


## Installation

If you have a working [Go installation](http://golang.org/doc/install), you
can easily get an up-to-date Grender binary.

    go get github.com/peterbourgon/grender


## Background

Grender is designed to power a large class of simple websites, not just blogs.
(Of course, Grender can power blogs, too.) It was created as a response to the
perceived complexity of Jekyll, both in terms of configuration and
implementation.


## Concepts

* **source directory**: specified by the commandline flag `-source` (default 
  `src`). Contains source files, which grender reads and processes.

* **target directory**: specified by the commandline flag `-target` (default 
  `tgt`). Every target file is rendered or copied from exactly one source file. 
  Not every source file produces a target file.

* **metadata**: JSON data, read from certain types of source files, and used to
  render target files. Associated with the source file that produced it.

* **the Stack**: a mapping of source filename to metadata. Provides "get"
  semantics, which composes the correct metadata for a target file. See
  [Metadata Composition](#Metadata-Composition), below, for more details.

* **global key**: a key present in every composed metadata, containing the
  specific composed metadata of every renderable source file. Specified by the
  commandline flag `-global.key` (default `files`). See
  [Global Key](#Global-Key), below, for more details.

* **template**: .md and .html files may contain template directives. Grender
  uses Go's [html/template](http://golang.org/pkg/html/template) engine. It's
  easy; see the examples directory for some use-cases.


## Metadata Composition

Grender's "secret sauce" is the way it builds up metadata. Basically, all
metadata is read and stored in the Stack. Then, every file that should be
rendered queries the Stack for the metadata that's relevant to it. Relevant
metadata is, in order of preference:

* metadata for the file itself
* metadata for the file's directory
* metadata for the file's directory's parent
* metadata for the file's directory's parent's parent
* ...and so on

"Closer" metadata is preferred to "further" metadata; that is, a key defined
in the file's metadata takes precedence over the same key defined in the
directory's metadata.

Declare metadata for a single file by putting a valid JSON object at the
beginning of the file, and separate it from the body with `---` + a newline.
Only .html and .md files may declare metadata this way.

Declare metadata for a directory (and all subdirectories) by putting a valid
JSON object in a .json file in that directory. Multiple .json files are
permitted; a single JSON object is composed from all of them.

That's all a bit abstract. Consider this example.

## Example

### foo/-.json

```
{"section": "foo"}
```

Every renderable file in the **foo** subdirectory, and all subdirectories, will
receive a metadata variable **section**, defined by default to be **foo**. Note
that we use the filename `-.json`, for lexical ordering reasons.

### foo/bar/index.html

```
{"subsection":"bar"}
---
<h1>Section {{ .section }}</h1>
<h2>subsection {{ .subsection }}</h2>
```

The **foo/bar/index.html** template inherits the **section** metadata, and also
defines a **subsection** metadata, accessible only to itself.


## Global Key

So far, we have enough in our toolbox to build a simple website. But more
dynamic features, like presenting list of blog entries, require a more
comprehensive view of the state of the website.

When grender processes source files and stores them into the Stack, it also
immediately queries the Stack for the composed metadata that's relevant to each
file. Grender pushes all of these composed metadatas into a separate piece of
metadata, tracked separately, called the Global Key.

When rendering a source file, grender queries the Stack for the specific,
composed metadata relevant to that file. But, it also injects the Global Key.
So, the template of the source file can read the composed metadata of any other
source file, or range over all the source files in a particular directory, or
anything else it wants to do.

Again, that's a bit abstract. Consider this example.

### blog/2001/01/01.md

```
{"title": "Alpha"}
---
...
```

### blog/2001/01/05.md

```
{"title": "Beta"}
---
...
```

### blog/index.html

```

{{ range .files.blog }}
	Blog posts from the year {{ . }}
	<ul>
		{{ range . }}
			<li><a href="{{ .url }}">{{ .title }}</a></li>
		{{ end }}
	</ul>
{{ end }}
```


## Autopopulated metadata

Every file gets some metadata for free. It can be overridden.

* TODO
* document
* those
* keys


## Usage

* TODO link to each example
