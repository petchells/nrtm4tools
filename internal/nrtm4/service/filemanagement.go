package service

import (
	"io"
	"log"
	"os"
	"strings"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

func initializeSource(repo persist.Repository, url string, notification nrtm4model.Notification, path string) error {
	var err error
	var file *os.File
	file, err = fileToDatabase(repo, url, notification.NrtmFile, persist.NotificationFile, path)
	if err != nil {
		return err
	}
	log.Println("DEBUG notification file.Name()", file.Name())

	log.Println("INFO file", file.Name())
	var snapshotFile *os.File
	snapshotFile, err = writeResourceToPath(notification.Snapshot.Url, path)
	// file, err = fileToDatabase(repo, notification.Snapshot.Url, nrtmFile, persist.SnapshotFile, path)
	log.Println("DEBUG wrote snapshotFile", snapshotFile.Name())
	return err
}

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
	ds := NrtmDataService{repo}
	return file, ds.saveState(url, nrtmFile, filetype, file)
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
