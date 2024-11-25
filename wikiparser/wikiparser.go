package wikiparser

import (
	"bufio"
	"compress/bzip2"
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/bodgit/sevenzip"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func extract_xml_to_page_by_golang(fileIo io.Reader, pageChan chan *Page, bar *progressbar.ProgressBar) {

	decoder := xml.NewDecoder(fileIo)
	for {
		token, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			} else {
				log.Fatal().Err(tokenErr).Msg("Error while decoding xml")
			}
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "page" {
				var page Page
				if err := decoder.DecodeElement(&page, &t); err != nil {
					log.Warn().Err(err).Msg("Error while decoding page")
				}
				pageChan <- &page
				bar.Add(1)
			}
		case xml.EndElement:

		}
	}
}

func extract_xml_to_page_by_xsp(fileIo io.Reader, pageChan chan *Page, bar *progressbar.ProgressBar) {

	br := bufio.NewReaderSize(fileIo, 65536)
	parser := xmlparser.NewXMLParser(br, "page")

	for pageXml := range parser.Stream() {
		if pageXml.Err != nil {
			log.Fatal().Err(pageXml.Err).Msg("Error while decoding xml")
		}

		ID, _ := strconv.ParseInt(pageXml.Childs["id"][0].InnerText, 10, 64)
		ns, _ := strconv.ParseInt(pageXml.Childs["ns"][0].InnerText, 10, 32)
		var page = Page{
			ID:    ID,
			Ns:    int32(ns),
			Title: pageXml.Childs["title"][0].InnerText,
		}

		if len(pageXml.Childs["redirect"]) > 0 {
			page.Redirect = &Redirect{
				Title: pageXml.Childs["redirect"][0].Attrs["title"],
			}
		}

		for _, revision := range pageXml.Childs["revision"] {
			revID, _ := strconv.ParseInt(revision.Childs["id"][0].InnerText, 10, 64)

			var parentid int64
			if len(revision.Childs["parentid"]) > 0 {
				parentid, _ = strconv.ParseInt(revision.Childs["parentid"][0].InnerText, 10, 64)
			}
			var comment string
			if len(revision.Childs["comment"]) > 0 {
				comment = revision.Childs["comment"][0].InnerText
			}
			Timestamp := revision.Childs["timestamp"][0].InnerText
			Model := revision.Childs["model"][0].InnerText
			Format := revision.Childs["format"][0].InnerText

			bytes, _ := strconv.ParseInt(revision.Childs["text"][0].Attrs["bytes"], 10, 32)

			var text = Text{
				Value:   revision.Childs["text"][0].InnerText,
				Bytes:   int32(bytes),
				Deleted: revision.Childs["text"][0].Attrs["deleted"],
			}
			ContID, _ := strconv.ParseInt(revision.Childs["contributor"][0].Attrs["id"], 10, 64)
			var contributor = Contributor{
				Username: revision.Childs["contributor"][0].Attrs["username"],
				Ip:       revision.Childs["contributor"][0].Attrs["ip"],
				ID:       ContID,
				Deleted:  revision.Childs["contributor"][0].Attrs["deleted"],
			}

			var revisionItem = Revision{
				ID:          revID,
				Text:        text,
				Parentid:    parentid,
				Timestamp:   Timestamp,
				Comment:     comment,
				Model:       Model,
				Format:      Format,
				Contributor: contributor,
			}
			page.Revisions = append(page.Revisions, revisionItem)
		}

		pageChan <- &page
		bar.Add(1)
	}

}

func ParseBzipXmlMixedFlow(filePathList []string, threadCount int) (<-chan *Page, <-chan string) {
	pageChan := make(chan *Page, 255)
	filePathChan := make(chan string, 64)
	completedFileChan := make(chan string, 2048)
	bar := progressbar.Default(-1)
	defer bar.Close()

	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
		close(filePathChan)
	}()

	wg := sync.WaitGroup{}
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {
		go func() {
			for filePath := range filePathChan {
				fileIo, err := os.Open(filePath)
				if err != nil {
					log.Fatal().Err(err).Str("filePath", filePath).Msg("Error while opening file")
				}
				bzio := bzip2.NewReader(fileIo)
				extract_xml_to_page_by_xsp(bzio, pageChan, bar)
				fileIo.Close()
				completedFileChan <- filePath
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pageChan)
		close(completedFileChan)
	}()

	return pageChan, completedFileChan
}

func Parse7zXmlMixedFlow(filePathList []string, threadCount int) (<-chan *Page, <-chan string) {
	pageChan := make(chan *Page, 255)
	filePathChan := make(chan string, 64)
	completedFileChan := make(chan string, 2048)

	bar := progressbar.Default(-1)
	defer bar.Close()

	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
		close(filePathChan)
	}()

	wg := sync.WaitGroup{}
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {

		go func() {
			for filePath := range filePathChan {
				r, err := sevenzip.OpenReader(filePath)
				if err != nil {
					log.Fatal().Err(err).Str("filePath", filePath).Msg("Error while opening file")
				}

				rc, err := r.File[0].Open()
				if err != nil {
					log.Fatal().Err(err).Msg("Error while opening 7z file")
				}

				extract_xml_to_page_by_xsp(rc, pageChan, bar)
				rc.Close()
				r.Close()
				completedFileChan <- filePath
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pageChan)
		close(completedFileChan)
	}()

	return pageChan, completedFileChan
}

func ParseBzipXmlSeparateFlow(filePathList []string, threadCount int, callback func(pageChan <-chan *Page, filePath string)) {

	filePathChan := make(chan string, 64)
	bar := progressbar.Default(-1)
	defer bar.Close()

	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
		close(filePathChan)
	}()

	wg := sync.WaitGroup{}
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {
		go func() {
			for filePath := range filePathChan {
				fileIo, err := os.Open(filePath)
				if err != nil {
					log.Fatal().Err(err).Str("filePath", filePath).Msg("Error while opening file")
				}
				bzio := bzip2.NewReader(fileIo)
				pageChan := make(chan *Page, 255)
				go func() {
					extract_xml_to_page_by_xsp(bzio, pageChan, bar)
					close(pageChan)
				}()
				callback(pageChan, filePath)
				fileIo.Close()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Parse7zXmlSeparateFlow(filePathList []string, threadCount int, callback func(pageChan <-chan *Page, filePath string)) {

	filePathChan := make(chan string, 2048)

	bar := progressbar.Default(-1)
	defer bar.Close()
	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
		close(filePathChan)
	}()

	wg := sync.WaitGroup{}
	wg.Add(threadCount)
	for i := 0; i < threadCount; i++ {

		go func() {
			for filePath := range filePathChan {
				r, err := sevenzip.OpenReader(filePath)
				if err != nil {
					log.Fatal().Err(err).Str("filePath", filePath).Msg("Error while opening file")
				}

				rc, err := r.File[0].Open()
				if err != nil {
					log.Fatal().Err(err).Msg("Error while opening 7z file")
				}
				pageChan := make(chan *Page, 255)
				go func() {
					extract_xml_to_page_by_xsp(rc, pageChan, bar)
					close(pageChan)
				}()
				callback(pageChan, filePath)
				rc.Close()
				r.Close()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
