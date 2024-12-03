package wikiparser

import (
	"testing"
)

func TestParse7zXML(t *testing.T) {
	pageChan, _ := Parse7zXmlMixedFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
		// "../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
	}, 10, []int32{0, 14})
	i := 0
	for range pageChan {
		i += 1
		// if i == 1 {
		// 	t.Logf("%+v", page)
		// }
	}
	if i != 562 {
		t.Error(i)
		t.Error("Not 562")
	}
}

func TestParseBzipXML(t *testing.T) {
	pageChan, _ := ParseBzipXmlMixedFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.bz2",
	}, 10, []int32{0, 14})

	i := 0
	for page := range pageChan {
		i += 1
		if i == 1 {
			t.Logf("%+v", page.Title)
		}
		// if page.Redirect != nil && page.Redirect.Title != "" {
		// 	t.Logf("%+v", page.Redirect.Title)
		// }
	}
}

func TestSpeparateSteam(t *testing.T) {
	Parse7zXmlSeparateFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
		// "/mnt/st01/wikipeida_download/20241101/enwiki-20241101-pages-meta-history10.xml-p5093021p5137508.7z",
	}, 10, []int32{0, 14}, func(pageChan <-chan *Page, filePath string) {
		i := 0
		for page := range pageChan {
			i += 1
			if i == 1 {
				t.Logf("%+v", page.Title)
			}
			// if page.Redirect != nil && page.Redirect.Title != "" {
			// 	t.Logf("%+v", page)
			// }
		}
		if i != 562 {
			t.Error(i)
			t.Error("Not 562")
		}
	})

	ParseBzipXmlSeparateFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.bz2",
	}, 10, []int32{0, 14}, func(pageChan <-chan *Page, filePath string) {
		i := 0
		for page := range pageChan {
			i += 1
			if i == 1 {
				t.Logf("%+v", page.Title)
			}
			// if page.Redirect != nil && page.Redirect.Title != "" {
			// 	t.Logf("%+v", page)
			// }
		}
		t.Log(filePath)
	})

}
