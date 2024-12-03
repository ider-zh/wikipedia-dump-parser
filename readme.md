# Wikipedia Dump Parser

This is a Go project that provides a parser for Wikipedia XML dumps. The parser can handle 7z and bzip2 compressed XML files.

## Features

- Parse 7z and bzip2 compressed Wikipedia XML dumps.
- Extract pages from the XML files.
- Output pages in a structured format.


## useage

`go get github.com/ider-zh/wikipedia-dump-parser`

```
impoer "github.com/ider-zh/wikipedia-dump-parser/wikiparser"

func main() {
	pageChan, _ := Parse7zXmlMixedFlow([]string{
		"./enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
	}, 10, []int32{0, 14})
	i := 0
	for range pageChan {
		i++
	}
    fmt.Println(i)
}
```

## test

 `go test -bench=. -benchmem -memprofile=mem.pprof -cpuprofile=cpu.pprof -blockprofile=block.pprof ./...`

 `go tool pprof -http="0.0.0.0:8080" mem.pprof  `

