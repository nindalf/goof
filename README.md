Goof
========

Goof (Go offer one file) is the Go version of [Woof](https://bitbucket.org/edu/woof/src/). It starts a server which serves a file.

Usage
----

To serve a file once on 127.0.0.1:8086

`goof -f <path-to-file>`

Additional options:

`goof -f <path-to-file> -c 2 -t 3600 -i 192.168.1.9 -p 3000`

* `f` - file path of the file that should be shared. Required.

* `c` - The number of times the file can be downloaded. Optional, default is 1. `c` = -1 indicates unlimited number of downloads.

* `t` - The time in minutes that the server runs before exiting. Optional, default is 0 (forever until parameter n is satisfied)

* `i` - The IP address on which the server should run. Optional, default is 127.0.0.1

* `p` - The port on which the server should listen. Optional, default is 8086