package gsnapshot

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (s GoogleSnapshotService) Create(diskId string, description string, zone string) (string, error) {
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
		SourceDisk:  diskId,
	}

	s.logger.Debug(googleSnapshotServiceLogTag, "Creating Google Snapshot with params: %#v", snapshot)
	operation, err := s.computeService.Disks.CreateSnapshot(s.project, gutil.ResourceSplitter(zone), diskId, snapshot).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Snapshot")
	}

	if _, err = s.operationService.Waiter(operation, zone, ""); err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Snapshot")
	}

	return snapshot.Name, nil
}
