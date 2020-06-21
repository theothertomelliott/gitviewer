# gitviewer

A web server to view the contents of a Git repository at any commit.

# Installation

A binary can be installed using `go get`:

    $ go get github.com/theothertomelliott/gitviewer/cmd/gitviewer

# Usage

The following command will start a server on port 8080, serving the contents of
the provided Git repo.

    $ gitviewer [repo path or url]

A Git repository may be specified as either an http(s) URL to a Git repo (as may be used with `git clone`), or a path to a local directory containing a Git repo.

Once started, you can browse to the server at: http://localhost:8080.

For a list of all config flags:

    $ gitviewer --help
