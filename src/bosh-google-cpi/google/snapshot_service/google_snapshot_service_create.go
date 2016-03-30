package snapshot

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
	"google.golang.org/api/compute/v1"
)

func (s GoogleSnapshotService) Create(diskID string, description string, zone string) (string, error) {
	uuidStr, err := s.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random Google Snapshot name")
	}

	if description == "" {
		description = googleSnapshotDescription
	}

	snapshot := &compute.Snapshot{
		Name:        fmt.Sprintf("%s-%s", googleSnapshotNamePrefix, uuidStr),
		Description: description,
	}

	s.logger.Debug(googleSnapshotServiceLogTag, "Creating Google Snapshot with params: %#v", snapshot)
	operation, err := s.computeService.Disks.CreateSnapshot(s.project, util.ResourceSplitter(zone), diskID, snapshot).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Snapshot")
	}

	if _, err = s.operationService.Waiter(operation, zone, ""); err != nil {
		s.cleanUp(snapshot.Name)
		return "", bosherr.WrapErrorf(err, "Failed to create Google Snapshot")
	}

	return snapshot.Name, nil
}

func (s GoogleSnapshotService) cleanUp(id string) {
	if err := s.Delete(id); err != nil {
		s.logger.Debug(googleSnapshotServiceLogTag, "Failed cleaning up Google Snapshot '%s': %#v", id, err)
	}
}
