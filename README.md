# Grender

Grender is a static site generator. It combines source files with metadata to
produce a website.

[![Build Status][1]][2]

[1]: https://secure.travis-ci.org/peterbourgon/grender.png
[2]: http://www.travis-ci.org/peterbourgon/grender


## Installation

If you have a working [Go installation](http://golang.org/doc/install), you
can easily get an up-to-date grender binary.

    go get github.com/peterbourgon/grender


## Background

Grender is designed to power a large class of simple websites, including but
not limited to blogs. It was created as a response to the perceived complexity
of Jekyll, both in terms of configuration and implementation.


## Usage

### Single file

Grender renders source files from the **source directory** (specified by the
commandline flag `-source`, default `src`) into the **target directory**
(`-target`, default `tgt`). Grender can render a single source file using
metadata provided in the file itself. Metadata is always valid JSON, and can be
put at the top of certain source files if it's separated by a line containing
only `---`.

See [the example][01].

[01]: http://github.com/peterbourgon/grender/blob/grender-2/examples/01-single-file


### Separate JSON

Metadata doesn't need to be in the source file directly. A valid .json file
provides metadata to every source file in the same directory. In case of
collision, grender prefers "closer" metadata; a file's specific metadata always
overrides a directory's metadata, for example. .json files are read in
lexigraphical order, before any source files are read.

See [the example][02]. Note that the .json file isn't copied to the target dir.

[02]: http://github.com/peterbourgon/grender/blob/grender-2/examples/02-separate-json


### Layering metadata

.json files provide metadata not only to every source file in the same
directory, but to all files in all subdirectories, too. In this way you can
"layer" metadata. The source file always receives a single, well-composed
metadata object for rendering.

[03]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata

See [the example][03]. The concept and application of composable metadata is
grender's Secret Sauce™.


### Imports

You can refer to the same content from multiple source files by using imports.
Save the shared content in a file; use the .source extension so grender knows
not to copy it to the target directory. Then use one of the import directives
to import it:

* `{{ importhtml "../relative/path.html.source" }}` for HTML snippets
* `{{ importcss "../relative/path.css.source" }}` for CSS snippets
* `{{ importjs "../relative/path.js.source" }}` for JS snippets

See [the example][04].

[04]: http://github.com/peterbourgon/grender/blob/grender-2/examples/04-imports


### Markdown and templates

Sometimes it's nice to specify a page merely as its content, and leave it to
the rendering engine to put that into a nice template. .md (Markdown) files
have this behavior by default. Grender expects to find a "template" key in
their metadata, and uses that filename as the template into which the rendered
Markdown is placed. Rendered content is available under the "content" key.

Template files should have the extension .template, so that grender knows not
to copy them to the target directory.

See [the example][05].

[05]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates

**Bonus**: if a Markdown filename matches the format YYYY-MM-DD-some-text.md, 
grender will treat that file as a "blog entry", and perform special behavior.
Given 2013-03-04-foo-bar-baz.md:

* default metadata key **title**, value "Foo bar baz"
* default metadata key **date**, value "2013 03 28"
* default target file is 2013/03/04/foo-bar-baz.html (relative to source)
* http-equiv refresh redirects to the target URL are written for all of the 
  following relative URLs:
** 2013/03/04/index.html
** 2013/03/4/index.html
** 2013/3/04/index.html
** 2013/3/4/index.html


### Discovering other files and metadata

So far we have enough tools to build a basic website. But we don't have any way
of linking to pages we don't explicitly know about. Grender solves this using
something called the Global Key (specified by the commandline flag
`-global.key`, default `files`). Any template may refer to any other file,
including all of its contextual metadata, via this key.

Say you want to build a list of every file available in the "blog" directory.
In your template, you can do:

```
<ul>
{{ range .files.blog }}
 <li> <a href="{{ .url }}">{{ .title }}</a> </li>
{{ end }}
</ul>
```

(`url` is a special key that grender autopopulates in the Global Key space for
every rendered file.) And what if you only want to list *most* of the files in
the "blog" directory? You can create `blog/default.json`:

```
{ "list": true }
```

And for the pages you don't want to include:

```
{ "list": false }
---
Content here
```

Then, in your template:

```
{{ range .files.blog }}
  {{ if .list }}
    <a href="{{ .url }}">{{ .title }}</a>
  {{ end }}
{{ end }}
```

See [the complete example][06].

[06]: http://github.com/peterbourgon/grender/blob/grender-2/examples/06-basic-blog


