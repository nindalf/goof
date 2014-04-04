Goof
========

Goof (Go offer one file) is the Go version of [Woof](https://bitbucket.org/edu/woof/src/). It is used to share a single file

Usage
----

To serve one file

`goof -f <path-to-file>`

`goof -f <path-to-file> -n 2 -t 3600`

* `f` - file path of the file that should be shared. Required.

* `n` - The number of times the file can be downloaded. Optional, default is 1.

* `t` - The time in minutes that the server runs before exiting. Optional, default is 0 (forever until parameter n is satisfied)