package parser

import (
	"compress/bzip2"
	"encoding/xml"
	"io"
	"os"

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

func ParseBzipXML(filePathList []string) chan *Page {
	pageChan := make(chan *Page, 1000)
	bar := progressbar.Default(-1)

	for _, filePath := range filePathList {
		go func() {
			fileIo, err := os.Open(filePath)
			if err != nil {
				log.Fatal().Err(err).Msg("Error while opening file")
			}
			defer fileIo.Close()

			bzio := bzip2.NewReader(fileIo)
			extract_xml_to_page(bzio, pageChan, bar)
		}()
	}
	return pageChan
}

func Parse7zXML(filePathList []string) chan *Page {
	pageChan := make(chan *Page, 1000)
	bar := progressbar.Default(-1)

	for _, filePath := range filePathList {
		go func() {
			r, err := sevenzip.OpenReader(filePath)
			if err != nil {
				log.Fatal().Err(err).Msg("Error while opening file")
			}
			defer r.Close()

			rc, err := r.File[0].Open()
			if err != nil {
				log.Fatal().Err(err).Msg("Error while opening 7z file")
			}
			defer rc.Close()

			extract_xml_to_page(rc, pageChan, bar)
		}()
	}
	return pageChan
}
