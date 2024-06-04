package service

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/jsonseq"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
)

// GZIPSnapshotExtension extension GZIP files
var GZIPSnapshotExtension = ".gz"

type fileManager struct {
	client Client
}

func (fm fileManager) ensureDirectoryExists(path string) error {
	var err error
	if _, err = os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return nil
	}
	err = os.Mkdir(path, 0755)
	if err != nil {
		return err
	}
	logger.Info("Created directory", "path", path)
	return nil
}

func (fm fileManager) fetchFileAndCheckHash(fileRef persist.FileRefJSON, path string) (*os.File, error) {
	var file *os.File
	_, err := os.Stat(filepath.Join(path, filepath.Base(fileRef.URL)))
	if os.IsNotExist(err) {
		logger.Info("Downloading file", "url", fileRef.URL)
		if file, err = fm.writeResourceToPath(fileRef.URL, path); err != nil {
			logger.Error("Failed to write file", "url", fileRef.URL, "path", path)
			return nil, err
		}
	} else {
		if file, err = os.Open(filepath.Join(path, filepath.Base(fileRef.URL))); err != nil {
			return nil, err
		}
	}
	sum, err := calcHash256(file)
	if err != nil {
		return nil, err
	}
	if sum != fileRef.Hash {
		if err = os.Rename(file.Name(), file.Name()+"-BADHASH"); err != nil {
			return nil, err
		}
		logger.Warn("Hash does not match the downloaded file. Try again", "file", file.Name(), "hash", fileRef.Hash, "calculated", sum)
		return nil, errors.New("hash does not match downloaded file")
	}
	return file, nil
}

func (fm fileManager) readJSONSeqRecords(
	file *os.File,
	fn jsonseq.RecordReaderFunc,
) error {

	var err error

	logger.Debug("opening for reading", "filename", file.Name())
	var reader io.Reader
	if reader, err = os.Open(file.Name()); err != nil {
		return err
	}
	var bufioReader *bufio.Reader
	if file.Name()[len(file.Name())-len(GZIPSnapshotExtension):] == GZIPSnapshotExtension {
		var gzreader *gzip.Reader
		if gzreader, err = gzip.NewReader(reader); err != nil {
			return err
		}
		bufioReader = bufio.NewReader(gzreader)
	} else {
		bufioReader = bufio.NewReader(reader)
	}
	err = jsonseq.ReadRecords(bufioReader, func(bytes []byte, err error) error {
		return fn(bytes, err)
	})
	return err
}

func (fm fileManager) writeResourceToPath(url string, path string) (*os.File, error) {
	fileName := filepath.Base(url)
	if f, err := os.Open(filepath.Join(path, fileName)); err == nil {
		return f, err
	}
	var reader io.Reader
	var err error
	if reader, err = fm.client.getResponseBody(url); err != nil {
		logger.Error("Failed to fetch file", url, err)
		return nil, err
	}
	return readerToFile(reader, path, fileName)
}

func (fm fileManager) downloadNotificationFile(url string) (persist.NotificationJSON, []error) {
	var notification persist.NotificationJSON
	var err error
	if notification, err = fm.client.getUpdateNotification(url); err != nil {
		logger.Error("fetching notificationFile", err)
		return notification, []error{err}
	}
	errs := fm.validateNotificationFile(notification)
	return notification, errs
}

func (fm fileManager) validateNotificationFile(file persist.NotificationJSON) []error {
	var errs []error
	if file.NrtmVersion != 4 {
		errs = append(errs, newNRTMServiceError("notificationFile nrtm version is not v4: '%v'", file.NrtmVersion))
	}
	if len(file.SessionID) < 36 {
		errs = append(errs, newNRTMServiceError("notificationFile session ID is not valid: '%v'", file.SessionID))
	}
	if len(file.Source) < 3 {
		errs = append(errs, newNRTMServiceError("notificationFile source is not valid: '%v'", file.Source))
	}
	if file.Version < 1 {
		errs = append(errs, newNRTMServiceError("notificationFile version must be positive: '%v'", file.NrtmVersion))
	}
	if len(file.SnapshotRef.URL) < 20 {
		errs = append(errs, newNRTMServiceError("notificationFile snapshot url is not valid: '%v'", file.SnapshotRef.URL))
	}
	return errs
}

func readerToFile(reader io.Reader, path string, fileName string) (*os.File, error) {
	var outFile *os.File
	var err error
	if outFile, err = os.Create(filepath.Join(path, fileName)); err != nil {
		logger.Error("Failed to open file on disk", err)
		return nil, err
	}
	// defer func() {
	// 	if err := outFile.Close(); err != nil {
	// 		panic(err)
	// 	}
	// }()
	if err = transferReaderToFile(reader, outFile); err != nil {
		logger.Error("writing file:", err)
		return nil, err
	}
	return outFile, err
}

func transferReaderToFile(from io.Reader, to *os.File) error {
	buf := make([]byte, fileWriteBufferLength)
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

func calcHash256(file *os.File) (string, error) {
	var err error
	if _, err = file.Seek(0, io.SeekStart); err != nil {
		logger.Error("Failed to seek(0) on downloaded file", "file", file.Name())
		return "", err
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	sum := hex.EncodeToString(hasher.Sum(nil))
	return sum, err
}
