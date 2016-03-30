package fakes

import (
	"bosh-google-cpi/google/snapshot_service"
)

type FakeSnapshotService struct {
	CreateCalled      bool
	CreateErr         error
	CreateID          string
	CreateDiskID      string
	CreateDescription string
	CreateZone        string

	DeleteCalled bool
	DeleteErr    error

	FindCalled   bool
	FindFound    bool
	FindSnapshot snapshot.Snapshot
	FindErr      error
}

func (s *FakeSnapshotService) Create(diskID string, description string, zone string) (string, error) {
	s.CreateCalled = true
	s.CreateDiskID = diskID
	s.CreateDescription = description
	s.CreateZone = zone
	return s.CreateID, s.CreateErr
}

func (s *FakeSnapshotService) Delete(id string) error {
	s.DeleteCalled = true
	return s.DeleteErr
}

func (s *FakeSnapshotService) Find(id string) (snapshot.Snapshot, bool, error) {
	s.FindCalled = true
	return s.FindSnapshot, s.FindFound, s.FindErr
}
