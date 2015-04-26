package ginstance

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"

	"github.com/frodenas/bosh-google-cpi/google/util"
)

func (i GoogleInstanceService) AttachedDisks(id string) (GoogleInstanceAttachedDisks, error) {
	i.logger.Debug(googleInstanceServiceLogTag, "Finding Google Disks attached to Google Instance '%s'", id)

	var disks GoogleInstanceAttachedDisks

	instance, found, err := i.Find(id, "")
	if err != nil {
		return disks, err
	}
	if !found {
		return nil, bosherr.Errorf("Google Instance '%s' not found", id)
	}

	for _, disk := range instance.Disks {
		if disk.Boot != true {
			disks = append(disks, gutil.ResourceSplitter(disk.Source))
		}
	}

	return disks, nil
}
