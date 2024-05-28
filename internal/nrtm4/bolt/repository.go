package bolt

import (
	"encoding/json"
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
	bolt "go.etcd.io/bbolt"
)

var stateBucketName = "state"

// BoltRepository uses BoltDB as a backing store
type BoltRepository struct {
	path string
	db   *bolt.DB
}

// Initialize implements the NRTM client Repository interface
func (r *BoltRepository) Initialize(path string) error {
	var err error
	r.path = path
	r.db, err = bolt.Open(r.path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(stateBucketName))
		return err
	})
}

// CreateSource implements the NRTM client Repository interface
func (r *BoltRepository) CreateSource(label string, source string, notificationURL string) (*persist.NRTMSource, error) {
	return &persist.NRTMSource{}, nil
}

func (r *BoltRepository) SaveFile(state *persist.NRTMFile) error {
	_, err := json.Marshal(*state)
	if err != nil {
		return err
	}
	// return r.db.Update(func(tx *bolt.Tx) error {
	// 	stateBucket := tx.Bucket([]byte(stateBucketName))
	// 	if err := stateBucket.Put([]byte(state.Source), stateBytes); err != nil {
	// 		return err
	// 	}
	// 	_, err := tx.CreateBucket([]byte(state.Source))
	// 	return err
	// })
	return nil
}

// SaveSnapshotObject implements the NRTM client Repository interface
func (r *BoltRepository) SaveSnapshotObject(source persist.NRTMSource, rpslObject rpsl.Rpsl) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		objectBucket := tx.Bucket([]byte(source.Source))
		return objectBucket.Put([]byte(rpslObject.ObjectType+" "+rpslObject.PrimaryKey), []byte(rpslObject.Payload))
	})
}

// SaveSnapshotObjects implements the NRTM client Repository interface
func (r *BoltRepository) SaveSnapshotObjects(source persist.NRTMSource, state persist.NRTMFile, rpslObjects []rpsl.Rpsl) error {
	for _, rpslObject := range rpslObjects {
		err := r.SaveSnapshotObject(source, rpslObject)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close implements the NRTM client Repository interface
func (r *BoltRepository) Close() error {
	return r.db.Close()
}
