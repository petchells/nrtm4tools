package service

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"slices"

	"github.com/petchells/nrtm4tools/internal/nrtm4/jsonseq"
	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
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

const numVersionsPerDirectory = 10000

// fetchFileAndCheckHash returns an open file pointer to the file in fileRef.URL
func (fm fileManager) fetchFileAndCheckHash(unfURL string, fileRef persist.FileRefJSON, basePath string) (*os.File, error) {
	fURL := fullURL(unfURL, fileRef.URL)
	if !validateURLString(fURL) {
		logger.Error("URL in fileRef cannot be parsed", "unfURL", unfURL, "fileRef.URL", fileRef.URL)
		return nil, errors.New("invalid URL in reference")
	}
	vdir := (fileRef.Version / numVersionsPerDirectory) * numVersionsPerDirectory
	subdir := filepath.Join(basePath, fmt.Sprintf("%d", vdir))
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		err = os.Mkdir(subdir, 0775)
		if err != nil {
			logger.Error("Failed to create subdirectory", "subdir", subdir, "error", err)
			return nil, err
		}
	}
	path := filepath.Join(subdir, filepath.Base(fURL))
	var file *os.File
	if file, err = os.Open(path); err != nil {
		UserLogger.Debug("Downloading file", "url", fURL, "path", path)
		if _, err = fm.writeResourceToPath(fURL, path); err != nil {
			logger.Error("Failed to write file", "url", fURL, "path", path)
			return nil, err
		}
		if file, err = os.Open(path); err != nil {
			logger.Error("Failed to open file", "url", fURL, "path", path)
			return nil, err
		}
	} else {
		UserLogger.Debug("Using existing file", "url", fURL, "path", path)
	}
	sum, err := calcHash256(file)
	if err != nil {
		return nil, err
	}
	if sum != fileRef.Hash {
		if err = os.Rename(file.Name(), file.Name()+"-BADHASH"); err != nil {
			return nil, err
		}
		UserLogger.Error("Hash does not match the downloaded file", "file", file.Name(), "hash", fileRef.Hash, "calculated", sum)
		return nil, ErrHashMismatch
	}
	UserLogger.Debug("File hash is ok", "file", file.Name())
	return file, nil
}

func (fm fileManager) readJSONSeqRecords(
	file *os.File,
	fn jsonseq.RecordReaderFunc,
) error {

	var err error

	logger.Info("Reading file", "filename", file.Name())
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

func (fm fileManager) writeResourceToPath(url string, fileName string) (*os.File, error) {
	if f, err := os.Open(fileName); err == nil {
		return f, err
	}
	var reader io.Reader
	var err error
	if reader, err = fm.client.getResponseBody(url); err != nil {
		logger.Error("Failed to fetch file", "url", url, "error", err)
		return nil, err
	}
	return readerToFile(reader, fileName)
}

func (fm fileManager) downloadNotificationFile(url string) (persist.NotificationJSON, error) {
	var notification persist.NotificationJSON
	var err error
	if notification, err = fm.client.getUpdateNotification(url); err != nil {
		logger.Error("getUpdateNotification returned an error", "error", err)
		return notification, err
	}
	return notification, validateNotificationFile(notification)
}

func validateNotificationFile(file persist.NotificationJSON) error {
	if file.NrtmVersion != 4 {
		return newNRTMServiceError("notificationFile nrtm version is not v4: '%v'", file.NrtmVersion)
	}
	if len(file.SessionID) < 36 {
		return newNRTMServiceError("notificationFile session ID is not valid: '%v'", file.SessionID)
	}
	if len(file.Source) < 1 {
		return newNRTMServiceError("notificationFile source name is not valid: '%v'", file.Source)
	}
	if file.Version < 1 {
		return newNRTMServiceError("notificationFile version must be positive: '%v'", file.NrtmVersion)
	}
	if len(file.SnapshotRef.URL) < 10 {
		return newNRTMServiceError("notificationFile snapshot url is not valid: '%v'", file.SnapshotRef.URL)
	}
	if len(file.DeltaRefs) == 0 {
		return ErrNRTM4NoDeltasInNotification
	}
	versions := make([]int64, len(file.DeltaRefs))
	for i, dr := range file.DeltaRefs {
		versions[i] = dr.Version
	}
	versionSet := util.NewSet(versions...)
	if len(versionSet) != len(versions) {
		logger.Error("Duplicate delta version found in notification file", "source", file.Source)
		return ErrNRTM4DuplicateDeltaVersion
	}
	slices.Sort(versions)
	lo := versions[0]
	hi := versions[len(versions)-1]
	if hi != file.Version {
		return ErrNRTM4NotificationVersionDoesNotMatchDelta
	}
	if lo+int64(len(versions)-1) != hi {
		logger.Error("Delta version is missing from the notification file", "source", file.Source)
		return ErrNRTM4NotificationDeltaSequenceBroken
	}
	return nil
}

func readerToFile(reader io.Reader, fileName string) (*os.File, error) {
	var outFile *os.File
	var err error
	if outFile, err = os.Create(fileName); err != nil {
		logger.Error("Failed to open file on disk", "error", err)
		return nil, err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			panic(err)
		}
	}()
	if err = transferReaderToFile(reader, outFile); err != nil {
		logger.Error("writing file:", "error", err)
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
