Goof
========

Goof (Go offer one file) is the Go version of [Woof](https://bitbucket.org/edu/woof/src/). It starts a server which serves a file or a folder. You can control the number of times the file is downloaded as well as the time the server is active.

Installing
---

If you have Go installed then you can run the following commands

```
go get github.com/nindalf/goof
//cd to directory
go install 
```

Or download from here for [Linux](https://github.com/nindalf/goof/releases/download/v0.9/goof) and [Mac](https://github.com/nindalf/goof/releases/download/v0.9/goof-mac)

Using
---

#####To serve a file or folder once

`goof -f /path/to/file`

Additional options:

`goof -f /path/to/file -c 2 -t 60 -i 192.168.1.9 -p 3000`

* `f` - file path of the file/folder that should be shared. Required.

* `c` - The number of times the file can be downloaded. Optional, default is 1. `c` = -1 indicates unlimited number of downloads. Flag is ignored if folder is served interactively.

* `t` - The time in minutes that the server runs before exiting. Optional, default is 0 (forever until parameter `c` is satisfied). 

* `i` - The IP address on which the server should run. Optional, default is all available addresses.

* `p` - The port on which the server should listen. Optional, default is 8086.

* `a` - Indicates if the folder to be served is to be archived. Optional, default is false (not archived and served interactively).

#####Uploads

`goof -u`

* `f` - Directory where the file should be saved. Optional, default is current working directory.

* `t`, `i`, and `p` are the same. `c` and `a` do not apply.
