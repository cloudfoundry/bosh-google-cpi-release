package machinetype

import (
	"bosh-google-cpi/util"
	"fmt"
	"strings"
)

func (m GoogleMachineTypeService) CustomLink(cpu int, ram int, zone string, machineSeries string) string {
	suffix := ""
	prefix := ""
	extendedThreshold := int(float64(cpu*1024) * 6.5)
	if ram > extendedThreshold {
		suffix = "-ext"
	}
	if machineSeries != "" && strings.ToLower(machineSeries) != "n1" {
		prefix = fmt.Sprintf("%s-", strings.ToLower(machineSeries))
	}
	return fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/machineTypes/%scustom-%d-%d%s", m.project, util.ResourceSplitter(zone), prefix, cpu, ram, suffix)
}
