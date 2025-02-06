package service

import (
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"

	"github.com/petchells/nrtm4client/internal/nrtm4/jsonseq"
	"github.com/petchells/nrtm4client/internal/nrtm4/persist"
	"github.com/petchells/nrtm4client/internal/nrtm4/rpsl"
	"github.com/petchells/nrtm4client/internal/nrtm4/util"
)

type rpslObjectParser struct{}

type rpslParserPool struct {
	Parsers chan rpslObjectParser
}

func newParserPool(limit int) *rpslParserPool {
	pool := rpslParserPool{}
	pool.Parsers = make(chan rpslObjectParser, limit)
	for range limit {
		pool.Parsers <- rpslObjectParser{}
	}
	return &pool
}

func (pool *rpslParserPool) Acquire() rpslObjectParser {
	return <-pool.Parsers
}

func (pool *rpslParserPool) Release(p rpslObjectParser) {
	pool.Parsers <- p
}

func (pool *rpslParserPool) Close() {
	close(pool.Parsers)
}

func (p *rpslObjectParser) bytesToRPSL(bytes []byte) *rpsl.Rpsl {
	so := new(persist.SnapshotObjectJSON)
	if err := json.Unmarshal(bytes, so); err != nil {
		logger.Warn("Failed to unmarshal RPSL string from", "so.Object", so.Object, "error", err)
		return nil
	}
	rpsl, err := rpsl.ParseFromJSONString(so.Object)
	if err != nil {
		logger.Warn("Failed to parse rpsl.Rpsl from", "so.Object", so.Object, "error", err)
	}
	return &rpsl
}

// CounterMsg is a message that can be sent to a counter
type CounterMsg int

const (
	// STOP tells the counter it will be closed
	STOP CounterMsg = iota
	// SUCCESS tells the counter to add one to the success count
	SUCCESS
	// FAILURE tells the counter to add one to the failure count
	FAILURE
	// REPORT tells the counter to print sth
	REPORT
)

func snapshotObjectInsertFunc(repo persist.Repository, source persist.NRTMSource, notification persist.NotificationJSON) jsonseq.RecordReaderFunc {

	var snapshotHeader *persist.SnapshotFileJSON
	var wg sync.WaitGroup

	objectList := util.NewLockingList[rpsl.Rpsl](rpslInsertBatchSize * 2)
	counterMsgChan := make(chan CounterMsg, 1000)
	expectHeader := true
	successCount := 0
	failureCount := 0

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				counterMsgChan <- REPORT
			case msg := <-counterMsgChan:
				switch msg {
				case SUCCESS:
					successCount++
				case FAILURE:
					failureCount++
				case REPORT:
					logger.Info("Parsing snapshot file", "objects", successCount, "failed", failureCount)
				case STOP:
					ticker.Stop()
					return
				}
			}
		}
	}()

	parserPool := newParserPool(4)
	incrementCounters := func(res *rpsl.Rpsl) {
		if obj := res; obj != nil {
			objectList.Add(*obj)
			counterMsgChan <- SUCCESS
		} else {
			counterMsgChan <- FAILURE
		}
		rpslObjects := objectList.GetBatch(rpslInsertBatchSize)
		if len(rpslObjects) > 0 {
			err := repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
			if err != nil {
				log.Fatalln("Error saving snapshot object", err)
			}
		}
	}

	return func(bytes []byte, err error) error {
		if err == io.EOF {
			// Expected error reading to end of snapshot objects
			parser := parserPool.Acquire()
			incrementCounters(parser.bytesToRPSL(bytes))
			parserPool.Release(parser)
			wg.Wait()
			parserPool.Close()
			counterMsgChan <- STOP
			close(counterMsgChan)
			logger.Info("Closed snapshot file", "numFailures", failureCount, "numSuccess", successCount)
			rpslObjects := objectList.GetAll()
			err = repo.SaveSnapshotObjects(source, rpslObjects, snapshotHeader.NrtmFileJSON)
			if err != nil {
				return err
			}
			source.Version = uint32(snapshotHeader.Version)
			_, err = repo.SaveSource(source, notification)
			return err
		} else if err != nil {
			logger.Warn("error reading jsonseq records.", "error", err)
			return err
		} else if expectHeader {
			// First record is the Snapshot header
			expectHeader = false
			sf := new(persist.SnapshotFileJSON)
			if err = json.Unmarshal(bytes, sf); err != nil {
				counterMsgChan <- FAILURE
				counterMsgChan <- STOP
				close(counterMsgChan)
				logger.Warn("error unmarshalling JSON. Expected SnapshotFile", "error", err)
				return err
			}
			if sf.Version != notification.SnapshotRef.Version {
				return ErrNRTM4FileVersionMismatch
			}
			snapshotHeader = sf
			counterMsgChan <- SUCCESS
			return nil
		} else {
			// Subsequent records are objects
			parser := parserPool.Acquire()
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer parserPool.Release(parser)
				incrementCounters(parser.bytesToRPSL(bytes))
			}()
			return nil
		}
	}
}
