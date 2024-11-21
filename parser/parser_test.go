package parser

import (
	"testing"
)

func TestParse7zXML(t *testing.T) {
	pageChan := Parse7zXML([]string{
		"/mnt/st01/wikipeida_download/20241101/enwiki-20241101-pages-meta-history19.xml-p29515305p29671963.7z",
		"/mnt/st01/wikipeida_download/20241101/enwiki-20241101-pages-meta-history27.xml-p68479621p68727588.7z",
	})
	i := 0
	for page := range pageChan {
		i += 1
		if i == 1 {
			t.Logf("%+v", page)
		}
		if i > 100 {
			break
		}
	}
}

func TestParseBzipXML(t *testing.T) {
	pageChan := ParseBzipXML([]string{
		"/home/ider/download/enwiki-20241001-pages-articles.xml.bz2",
	})

	i := 0
	for page := range pageChan {
		i += 1
		if i == 1 {
			t.Logf("%+v", page)
		}
		if i > 100 {
			break
		}
	}
}
