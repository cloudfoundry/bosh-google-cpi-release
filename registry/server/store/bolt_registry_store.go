package store

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	"github.com/boltdb/bolt"
)

const BoltRegistryStoreLogTag = "BoltRegistryStore"
const BoltRegistryStoreFileMode = 0600
const BoltRegistryStoreFileLockTimeout = 1
const BoltRegistryStoreBucketName = "Registry"

type BoltRegistryStore struct {
	config BoltRegistryStoreConfig
	logger boshlog.Logger
}

func NewBoltRegistryStore(
	config BoltRegistryStoreConfig,
	logger boshlog.Logger,
) BoltRegistryStore {
	return BoltRegistryStore{
		config: config,
		logger: logger,
	}
}

func (s BoltRegistryStore) openDB() (db *bolt.DB, err error) {
	dbOptions := &bolt.Options{
		Timeout: BoltRegistryStoreFileLockTimeout * time.Second,
	}
	db, err = bolt.Open(s.config.DBFile, BoltRegistryStoreFileMode, dbOptions)
	if err != nil {
		return db, bosherr.WrapError(err, "Opening Bolt database")
	}

	return db, nil
}

func (s BoltRegistryStore) Delete(key string) error {
	db, err := s.openDB()
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting key '%s'", key)
	}
	defer db.Close()

	s.logger.Debug(BoltRegistryStoreLogTag, "Deleting key '%s'", key)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BoltRegistryStoreBucketName))
		if bucket != nil {
			return bucket.Delete([]byte(key))
		}
		return nil
	})
	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting key '%s'", key)
	}

	return nil
}

func (s BoltRegistryStore) Get(key string) (string, bool, error) {
	db, err := s.openDB()
	if err != nil {
		return "", false, bosherr.WrapErrorf(err, "Reading key '%s'", key)
	}
	defer db.Close()

	var value []byte
	s.logger.Debug(BoltRegistryStoreLogTag, "Reading key '%s'", key)
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BoltRegistryStoreBucketName))
		if bucket != nil {
			value = bucket.Get([]byte(key))
		}
		return nil
	})
	if value != nil {
		return string(value), true, nil
	}

	return "", false, nil
}

func (s BoltRegistryStore) Save(key string, value string) error {
	db, err := s.openDB()
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving key '%s'", key)
	}
	defer db.Close()

	s.logger.Debug(BoltRegistryStoreLogTag, "Saving key '%s'", key)
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BoltRegistryStoreBucketName))
		if err != nil {
			return bosherr.WrapErrorf(err, "Creating bucket '%s'", BoltRegistryStoreBucketName)
		}
		return bucket.Put([]byte(key), []byte(value))
	})
	if err != nil {
		return bosherr.WrapErrorf(err, "Saving key '%s'", key)
	}

	return nil
}
