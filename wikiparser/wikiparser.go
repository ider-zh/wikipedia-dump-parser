package wikiparser

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
				extract_xml_to_page(bzio, pageChan, bar)
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

				extract_xml_to_page(rc, pageChan, bar)
				rc.Close()
				r.Close()
				completedFileChan <- filePath
			}
		}()
	}

	go func() {
		wg.Wait()
		close(pageChan)
		close(completedFileChan)
	}()

	return pageChan, completedFileChan
}

type CallbackFunc func(pageChan <-chan *Page, filePath string)

func ParseBzipXmlSeparateFlow(filePathList []string, threadCount int, callback CallbackFunc) {

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
					extract_xml_to_page(bzio, pageChan, bar)
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

func Parse7zXmlSeparateFlow(filePathList []string, threadCount int, callback CallbackFunc) {

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
					extract_xml_to_page(rc, pageChan, bar)
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
