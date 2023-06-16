package service

import (
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func fileToDatabase(repo persist.Repository, url string, nrtmFile nrtm4model.NrtmFile, filetype persist.NTRMFileType, path string) (*os.File, error) {
	var file *os.File
	var err error
	if file, err = writeResourceToPath(url, path); err != nil {
		return file, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	state := persist.NRTMState{
		ID:       0,
		Created:  time.Now(),
		Source:   nrtmFile.Source,
		Version:  nrtmFile.Version,
		URL:      url,
		Type:     filetype,
		FileName: file.Name(),
	}
	err = repo.SaveState(state)
	return file, err
}

func writeResourceToPath(url string, path string) (*os.File, error) {
	fileName := url[strings.LastIndex(url, "/")+1:]
	var reader io.ReadCloser
	var httpClient HttpClient
	var outFile *os.File
	var err error
	if reader, err = httpClient.fetchFile(url); err != nil {
		log.Println("ERROR Failed to fetch file", url, err)
		return nil, err
	}
	if outFile, err = os.Create(path + "/" + fileName); err != nil {
		log.Println("ERROR Failed to open file on disk", err)
		return nil, err
	}
	if err = transferReaderToFile(reader, outFile); err != nil {
		log.Println("ERROR writing file")
		return nil, err
	}
	return outFile, err
}

func transferReaderToFile(from io.ReadCloser, to *os.File) error {
	buf := make([]byte, FILE_BUFFER_LENGTH)
	for {
		n, err := from.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := to.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
