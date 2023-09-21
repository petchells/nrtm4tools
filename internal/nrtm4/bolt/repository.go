package bolt

import (
	"encoding/json"
	"time"

	"gitlab.com/etchells/nrtm4client/internal/nrtm4/nrtm4model"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/persist"
	"gitlab.com/etchells/nrtm4client/internal/nrtm4/rpsl"
	bolt "go.etcd.io/bbolt"
)

var stateBucketName = "state"

type BoltRepository struct {
	path string
	db   *bolt.DB
}

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

func (r *BoltRepository) GetState(source string) (persist.NRTMState, error) {
	state := persist.NRTMState{}
	err := r.db.View(func(tx *bolt.Tx) error {
		stateBucket := tx.Bucket([]byte(stateBucketName))
		if stateBucket == nil {
			return &persist.ErrStateNotInitialized
		}
		stateValue := stateBucket.Get([]byte(source))
		if stateValue == nil {
			return &persist.ErrStateNotInitialized
		}
		return json.Unmarshal(stateValue, &state)
	})
	return state, err
}

func (r *BoltRepository) SaveState(state *persist.NRTMState) error {
	stateBytes, err := json.Marshal(*state)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		stateBucket := tx.Bucket([]byte(stateBucketName))
		if err := stateBucket.Put([]byte(state.Source), stateBytes); err != nil {
			return err
		}
		_, err := tx.CreateBucket([]byte(state.Source))
		return err
	})
}

func (r *BoltRepository) SaveSnapshotFile(persist.NRTMState, nrtm4model.SnapshotFile) error {
	return nil
}

func (r *BoltRepository) SaveSnapshotObject(state persist.NRTMState, rpslObject rpsl.Rpsl) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		objectBucket := tx.Bucket([]byte(state.Source))
		return objectBucket.Put([]byte(rpslObject.ObjectType+" "+rpslObject.PrimaryKey), []byte(rpslObject.Payload))
	})
}

func (r *BoltRepository) Close() error {
	return r.db.Close()
}
