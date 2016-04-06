package machinetype

import (
	"fmt"

	"bosh-google-cpi/util"
)

func (m GoogleMachineTypeService) CustomLink(cpu int, ram int, zone string) string {
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes/custom-%d-%d", m.project, util.ResourceSplitter(zone), cpu, ram)
}
