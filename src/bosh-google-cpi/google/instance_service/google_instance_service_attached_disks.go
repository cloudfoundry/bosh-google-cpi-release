package instance

import (
	"bosh-google-cpi/api"
	"bosh-google-cpi/util"
)

func (i GoogleInstanceService) AttachedDisks(id string) (AttachedDisks, error) {
	i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Disks attached to Google Instance '%s'", id)

	var disks AttachedDisks

	instance, found, err := i.Find(id, "")
	if err != nil {
		return disks, err
	}
	if !found {
		return disks, api.NewVMNotFoundError(id)
	}

	for _, disk := range instance.Disks {
		if disk.Boot != true {
			disks = append(disks, util.ResourceSplitter(disk.Source))
		}
	}

	return disks, nil
}
