# Kognit

Kognit is a tool developed for the sake of self development regarding file compression and the Go programming languuage.

## Features

- [X] Directory compression
  - [X] Zip
  - [X] Tar with gzip
- [ ] File compression
  - [ ] Flate
  - [ ] Deflate
  - [ ] gzip
  - [ ] Huffman
  - [ ] LZW
  - [ ] RLE
- [ ] Image compression
  - [ ] JPEG
  - [ ] JPEG2000
  - [ ] PNG
  - [ ] GIF
- [ ] Sound compression

## Usage

For now this tool can be used as an importable package with already some of its features implemented

To compress a directory's contents you can do:

``` Go
package main

import (
    "github.com/csothen/kognit"
)

func main(){
    // Will encode the directory into a zip file
    if err := kognit.ZIP.Encode("myDir"); err != nil {
        panic(err)
    }
}
```

All future features should follow the same usage flow:

- Define what type of compression you wish to perform by selecting it from the kognit package
- Call its Encode method to compress the data
- Call its Decode method when you want to decompress the data

## Future development

For the tools future it is intended that it can be used as:

- An importable package (in progress)
- A command line interface (to do)
- The backend for a web application (to do)
