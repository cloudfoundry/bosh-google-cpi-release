package disk

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"bosh-google-cpi/util"
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
	operation, err := d.computeService.Disks.Insert(d.project, util.ResourceSplitter(zone), disk).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk")
	}

	if _, err = d.operationService.Waiter(operation, zone, ""); err != nil {
		d.cleanUp(disk.Name)
		return "", bosherr.WrapErrorf(err, "Failed to create Google Disk")
	}

	return disk.Name, nil
}

func (d GoogleDiskService) cleanUp(id string) {
	if err := d.Delete(id); err != nil {
		d.logger.Debug(googleDiskServiceLogTag, "Failed cleaning up Google Disk '%s': %#v", id, err)
	}
}
