package wikiparser

import (
	"testing"
)

func TestParse7zXML(t *testing.T) {
	pageChan, _ := Parse7zXmlMixedFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
		// "../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
	}, 10)
	i := 0
	for range pageChan {
		i += 1
		// if i == 1 {
		// 	t.Logf("%+v", page)
		// }
	}
	if i != 1803 {
		t.Error(i)
		t.Error("Not 1803")
	}
}

func TestParseBzipXML(t *testing.T) {
	pageChan, _ := ParseBzipXmlMixedFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.bz2",
	}, 10)

	i := 0
	for page := range pageChan {
		i += 1
		if i == 1 {
			t.Logf("%+v", page)
		}
		// if page.Redirect != nil && page.Redirect.Title != "" {
		// 	t.Logf("%+v", page.Redirect.Title)
		// }
	}
}

func TestSpeparateSteam(t *testing.T) {
	Parse7zXmlSeparateFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
		// "../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.7z",
	}, 10, func(pageChan <-chan *Page, filePath string) {
		i := 0
		for page := range pageChan {
			i += 1
			if i == 1 {
				t.Logf("%+v", page)
			}
			// if page.Redirect != nil && page.Redirect.Title != "" {
			// 	t.Logf("%+v", page)
			// }
		}
		if i != 1803 {
			t.Error(i)
			t.Error("Not 1803")
		}
	})

	ParseBzipXmlSeparateFlow([]string{
		"../testdata/enwiki-20241101-pages-meta-history23.xml-p50562218p50564553.bz2",
	}, 10, func(pageChan <-chan *Page, filePath string) {
		i := 0
		for page := range pageChan {
			i += 1
			if i == 1 {
				t.Logf("%+v", page)
			}
			// if page.Redirect != nil && page.Redirect.Title != "" {
			// 	t.Logf("%+v", page)
			// }
		}
		t.Log(filePath)
	})

}
