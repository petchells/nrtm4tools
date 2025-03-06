package bolt

import (
	"errors"
	"time"

	"github.com/petchells/nrtm4tools/internal/nrtm4/persist"
	"github.com/petchells/nrtm4tools/internal/nrtm4/rpsl"
	bolt "go.etcd.io/bbolt"
)

var stateBucketName = "state"

// ErrNotImplemented but it will be "real soon"
var ErrNotImplemented = errors.New("function not implemented for BoltDB")

// BBolt uses BoltDB as a backing store
type BBolt struct {
	db *bolt.DB
}

// Initialize implements the NRTM client Repository interface
func (r *BBolt) Initialize(path string) error {
	var err error
	r.db, err = bolt.Open(path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(stateBucketName))
		return err
	})
}

// SaveSource implements the NRTM client Repository interface
func (r *BBolt) SaveSource(source persist.NRTMSource, notificationURL string) (persist.NRTMSource, error) {
	return persist.NRTMSource{}, nil
}

// ListSources TODO: implement
func (r *BBolt) ListSources() ([]persist.NRTMSource, error) {
	return []persist.NRTMSource{}, nil
}

// SaveSnapshotObject implements the NRTM client Repository interface
func (r *BBolt) SaveSnapshotObject(source persist.NRTMSource, rpslObject rpsl.Rpsl) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		objectBucket := tx.Bucket([]byte(source.Source))
		return objectBucket.Put([]byte(rpslObject.ObjectType+" "+rpslObject.PrimaryKey), []byte(rpslObject.Payload))
	})
}

// SaveSnapshotObjects implements the NRTM client Repository interface
func (r *BBolt) SaveSnapshotObjects(source persist.NRTMSource, state persist.NRTMFile, rpslObjects []rpsl.Rpsl) error {
	for _, rpslObject := range rpslObjects {
		err := r.SaveSnapshotObject(source, rpslObject)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close implements the NRTM client Repository interface
func (r *BBolt) Close() error {
	return r.db.Close()
}
