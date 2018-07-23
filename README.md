# goform

Various form handler for go

## Why ?

1. Any file is a form of data, is readable and writable if you follow some rules.
2. An editor is a viewer of some data structure, maybe we can just access the underlying data structure with goform.  
3. Json is a good enough uniform of data structure, maybe we can edit any form in json and then convert back.

## Quit start

```bash
# Install
go get github.com/wenerme/goform/cmd/goform/...

# Convert met to json
goform convert -i server.met -o server.json -O indent
# Edit
nano server.json
# Convert back
goform convert -i server.json -o server.met
```

## Supported extension
* `.met`
  * [server.met](http://wiki.amule.org/wiki/Server.met_file) for ed2k
* `.json`
