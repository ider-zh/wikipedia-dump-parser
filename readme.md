# Wikipedia Dump Parser

This is a Go project that provides a parser for Wikipedia XML dumps. The parser can handle 7z and bzip2 compressed XML files.

## Features

- Parse 7z and bzip2 compressed Wikipedia XML dumps.
- Extract pages from the XML files.
- Output pages in a structured format.


## test

 go test -bench=. -benchmem -memprofile=mem.pprof -cpuprofile=cpu.pprof -blockprofile=block.pprof ./...

 go tool pprof -http="0.0.0.0:8080" mem.pprof  

 