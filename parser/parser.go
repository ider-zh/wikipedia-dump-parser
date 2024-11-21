package parser

import (
	"compress/bzip2"
	"encoding/xml"
	"io"
	"os"
	"sync"

	"github.com/bodgit/sevenzip"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/schollz/progressbar/v3"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func extract_xml_to_page(fileIo io.Reader, pageChan chan *Page, bar *progressbar.ProgressBar) {
	decoder := xml.NewDecoder(fileIo)
	for {
		t, tokenErr := decoder.Token()
		if tokenErr != nil {
			if tokenErr == io.EOF {
				break
			} else {
				log.Fatal().Err(tokenErr).Msg("Error while decoding xml")
			}
		}
		switch t := t.(type) {
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

func ParseBzipXML(filePathList []string, threadCount int) chan *Page {
	pageChan := make(chan *Page, 1000)
	bar := progressbar.Default(-1)
	filePathChan := make(chan string, 1000)

	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
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
				extract_xml_to_page(bzio, pageChan, bar)
				fileIo.Close()
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pageChan)
	}()

	return pageChan
}

func Parse7zXML(filePathList []string, threadCount int) chan *Page {
	pageChan := make(chan *Page, 1000)
	bar := progressbar.Default(-1)

	filePathChan := make(chan string, 1000)

	go func() {
		for _, filePath := range filePathList {
			filePathChan <- filePath
		}
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

				extract_xml_to_page(rc, pageChan, bar)
				rc.Close()
				r.Close()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(pageChan)
	}()

	return pageChan
}
