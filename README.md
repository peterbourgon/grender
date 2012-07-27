# grender

A different take on a static site generator.

# Desires

Grender is designed to manage a particular kind of website. A grenderable
(grendered?) site will

* Probably have a blog section, where identially-styled "posts" or "articles"
  are kept and organized by date.

* Probably have many arbitrary content pages, which will have unique designs,
  and be referenced by name (by path).

* Definitely have a front page, probably with its own unique design, and
  probably taking dynamic elements from other parts of the site.

Using grender day-to-day should be something like this: you keep your website
data in "source" form -- most content in Markdown, maybe some pages in raw
HTML. That data is probably stored in a Github repository (or equivalent). To
update your website, you add new source content, and run the `grender`
executable, probably without any arguments. An "output" subdirectory is created
and filled with the "compiled" version of your website.

# Design concepts

There are two types of input files. Source files contain content, plus limited
metadata. Template files contain markup to render source. 

## Source files

Source files begin with a metadata section, which must end with the
`--metadata-delimiter`. They then continue with plaintext content, likely in
Markdown format.

Source file metadata is parsed as YAML. Some predefined keys are used (and
consumed) to control grender's behavior:

* `template` - template file used to render this source file

All other keys are passed directly to the templating engine.

The name of the output file, and its location in the output hierarchy, is
determined by rules. Those rules are as follows:

* If a source file's basename matches the pattern `YYYY-MM-DD-title-text`, it's
  considered a blog entry. A blog entry is rendered in the `--blog`
  subdirectory of the `--output` path, with no directory hierarchy.
  
* Otherwise, a source file is considered arbitrary content. It's rendered in
  the `--output` path with the same relative hierarchy as it has in the
  `--source` path. 

All output files are named as their source file's basename, plus the
`--output-extension`. For example, `mypage.txt` becomes `mypage.html`.

### Blog entries

A "blog archive" page may wish to provide links to all blog entries in a
particular date range. And the front-page may wish to display the N most recent
blog entries directly inline. For this reason grender must treat blog entries
in a special way.

### Arbitrary content

TODO

## Template files

TODO

# Errata

Grender's behavior will be controlled exclusively with commandline flags. It
will not support configuration files or pay attention to environment variables.
Nonstandard grender use cases should be handled by writing a wrapper script.

Grender should deduce as much of its behavior as possible from implicit
metadata: [source] file name and location. Convention over configuration.

Nontrivial URL mapping schemes should be managed in the web server
configuration. Grender works best when you're able to manipulate that config.

