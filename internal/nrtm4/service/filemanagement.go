package service

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jmank88/jsonseq"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

type fileManager struct {
	repo   persist.Repository
	client Client
}

func (fm fileManager) initializeSource(url string, path string, notification nrtm4model.Notification) error {

	var err error
	// var file *os.File
	// file, err = fm.fileToDatabase(url, notification.NrtmFile, persist.NotificationFile, path)
	// if err != nil {
	// 	return err
	// }
	// log.Println("DEBUG notification file.Name()", file.Name())

	// log.Println("INFO file", file.Name())
	var snapshotFile *os.File
	if snapshotFile, err = fm.writeResourceToPath(notification.Snapshot.Url, path); err != nil {
		return err
	}
	log.Println("DEBUG wrote snapshotFile", snapshotFile.Name())
	// TODO: parse file
	var reader io.Reader
	if reader, err = os.Open(snapshotFile.Name()); err != nil {
		return err
	}
	var gzreader *gzip.Reader
	if gzreader, err = gzip.NewReader(reader); err != nil {
		return err
	}
	d := jsonseq.NewDecoder(gzreader)
	for {
		var i interface{}
		if err := d.Decode(&i); err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
		} else {
			fmt.Println(i)
		}
	}
	// file, err = fileToDatabase(repo, notification.Snapshot.Url, nrtmFile, persist.SnapshotFile, path)
	return err
}

func (fm fileManager) fileToDatabase(url string, nrtmFile nrtm4model.NrtmFile, filetype persist.NTRMFileType, path string) (*os.File, error) {
	var file *os.File
	var err error
	if file, err = fm.writeResourceToPath(url, path); err != nil {
		return file, err
	}
	// defer func() {
	// 	if err := file.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()
	ds := NrtmDataService{fm.repo}
	return file, ds.saveState(url, nrtmFile, filetype, file)
}

func (fm fileManager) writeResourceToPath(url string, path string) (*os.File, error) {
	fileName := filepath.Base(url)
	var reader io.Reader
	var err error
	if reader, err = fm.client.fetchFile(url); err != nil {
		log.Println("ERROR Failed to fetch file", url, err)
		return nil, err
	}
	return readerToFile(reader, path, fileName)
}

func readerToFile(reader io.Reader, path string, fileName string) (*os.File, error) {
	var outFile *os.File
	var err error
	if outFile, err = os.Create(filepath.Join(path, fileName)); err != nil {
		log.Println("ERROR Failed to open file on disk", err)
		return nil, err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			panic(err)
		}
	}()
	if err = transferReaderToFile(reader, outFile); err != nil {
		log.Println("ERROR writing file")
		return nil, err
	}
	return outFile, err
}

func transferReaderToFile(from io.Reader, to *os.File) error {
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
