package machinetype

import (
	"fmt"

	"bosh-google-cpi/util"
)

func (m GoogleMachineTypeService) CustomLink(cpu int, ram int, zone string) string {
	suffix := ""
	extendedThreshold := int(float64(cpu*1024) * 6.5)
	if ram > extendedThreshold {
		suffix = "-ext"
	}
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes/custom-%d-%d%s", m.project, util.ResourceSplitter(zone), cpu, ram, suffix)
}
