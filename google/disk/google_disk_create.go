package gdisk

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
	"google.golang.org/api/compute/v1"
)

func (d GoogleDiskService) Create(size int, diskType string, zone string) (string, error) {
	uuidStr, err := d.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random Google Disk name")
	}

	disk := &compute.Disk{
		Name:        fmt.Sprintf("%s-%s", googleDiskNamePrefix, uuidStr),
		Description: googleDiskDescription,
		SizeGb:      int64(size),
	}

	if diskType != "" {
		disk.Type = diskType
	}

	d.logger.Debug(googleDiskServiceLogTag, "Creating Google Disk with params: %#v", disk)
	operation, err := d.computeService.Disks.Insert(d.project, gutil.ResourceSplitter(zone), disk).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk")
	}

	if _, err = d.operationService.Waiter(operation, zone, ""); err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk")
	}

	return disk.Name, nil
}
