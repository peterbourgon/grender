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

See [sample source][01-src] and [rendered output][01-tgt].

[01-src]: http://github.com/peterbourgon/grender/blob/grender-2/examples/01-single-file/src/index.html
[01-tgt]: http://github.com/peterbourgon/grender/blob/grender-2/examples/01-single-file/tgt/index.html


### Separate JSON

Metadata doesn't need to be in the source file directly. A valid .json file
provides metadata to every source file in the same directory. In case of
collision, grender prefers "closer" metadata; a file's specific metadata always
overrides a directory's metadata, for example.

See [sample .json file][02-src-json], [sample .html file][02-src-html] and
[rendered output][02-tgt-html]. Note that the .json file is not copied to the
[target directory][02-tgt-dir].

[02-src-json]: http://github.com/peterbourgon/grender/blob/grender-2/examples/02-separate-json/src/-.json
[02-src-html]: http://github.com/peterbourgon/grender/blob/grender-2/examples/02-separate-json/src/index.html
[02-tgt-html]: http://github.com/peterbourgon/grender/blob/grender-2/examples/02-separate-json/tgt/index.html
[02-tgt-dir]: http://github.com/peterbourgon/grender/blob/grender-2/examples/02-separate-json/tgt


### Layering metadata

.json files provide metadata not only to every source file in the same
directory, but to all files in all subdirectories, too. In this way you can
"layer" metadata. The source file always receives a single, well-composed
metadata object for rendering.

Given a [document][03-src-document], [section][03-src-section],
[subsection][03-src-subsection] hierarchy, a [source file][03-src-foo] produces
this [rendered output][03-tgt].

[03-src-document]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata/src/-.json
[03-src-section]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata/src/1/-.json
[03-src-subsection]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata/src/1/foo/-.json
[03-src-foo]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata/src/1/foo/index.html
[03-tgt]: http://github.com/peterbourgon/grender/blob/grender-2/examples/03-layering-metadata/tgt/1/foo/index.html

The concept and application of composable metadata is grender's Secret Sauce™.


### Imports

You can refer to the same content from multiple source files by using imports.
Save the shared content in a file; use the .source extension so grender knows
not to copy it to the target directory. Then use one of the import directives
to import it:

* `{{ importhtml "../relative/path.html.source" }}` for HTML snippets
* `{{ importcss "../relative/path.css.source" }}` for CSS snippets
* `{{ importjs "../relative/path.js.source" }}` for JS snippets

See it in action: [source document][04-src], [CSS .source][04-css], 
[HTML .source][04-html], and the [rendered output][04-tgt].

[04-src]: http://github.com/peterbourgon/grender/blob/grender-2/examples/04-imports/src/index.html
[04-css]: http://github.com/peterbourgon/grender/blob/grender-2/examples/04-imports/src/my.css.source
[04-html]: http://github.com/peterbourgon/grender/blob/grender-2/examples/04-imports/src/my.html.source
[04-tgt]: http://github.com/peterbourgon/grender/blob/grender-2/examples/04-imports/tgt/index.html


### Markdown and templates

Sometimes it's nice to specify a page merely as its content, and leave it to
the rendering engine to put that into a nice template. .md (Markdown) files
have this behavior by default. Grender expects to find a "template" key in
their metadata, and uses that filename as the template into which the rendered
Markdown is placed. Rendered content is available under the "content" key.

Template files should have the extension .template, so that grender knows not
to copy them to the target directory.

In action: directory-wide [.json file][05-src-json] specifies a
[template][05-src-template]; individual .md files ([one][05-src-01],
[two][05-src-02]) specify titles, only; [rendered output][05-tgt].

[05-src-json]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates/src/-.json
[05-src-template]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates/src/entry.template
[05-src-01]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates/src/01.md
[05-src-02]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates/src/02.md
[05-tgt]: http://github.com/peterbourgon/grender/blob/grender-2/examples/05-templates/tgt


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
every rendered file.) See [this example][06] for a use-case.

[06]: http://github.com/peterbourgon/grener/blob/grender-2/examples/06-basic-blog


