package action

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type DesiredVMSpec struct {
	CPU               int `json:"cpu"`
	RAM               int `json:"ram"`
	EphemeralDiskSize int `json:"ephemeral_disk_size"`
}

type CalculateVMCloudProperties struct{}

func NewCalculateVMCloudProperties() CalculateVMCloudProperties {
	return CalculateVMCloudProperties{}
}

const (
	gb                = 1024
	minMemoryPerCPU   = 0.9 * gb
	memoryGranularity = 256
)

const (
	NoCPUErr       = "CPU must be greater than 0"
	RamPerCPUErr   = "RAM per CPU must be at least 922 MB"
	RamMultipleErr = "RAM must be specified as a multiple of 256 MB"
)

func (CalculateVMCloudProperties) Run(desired DesiredVMSpec) (VMCloudProperties, error) {
	var errs []string
	if desired.CPU <= 0 {
		errs = append(errs, NoCPUErr)
	}
	minRam := float64(desired.CPU) * float64(minMemoryPerCPU)
	if float64(desired.RAM) < minRam {
		errs = append(errs, RamPerCPUErr)
	}
	if desired.RAM%memoryGranularity != 0 {
		errs = append(errs, RamMultipleErr)
	}

	if len(errs) != 0 {
		return VMCloudProperties{}, bosherr.Errorf("invalid desired instance specification: %s", strings.Join(errs, ", "))
	}

	return VMCloudProperties{
		CPU:            desired.CPU,
		RAM:            desired.RAM,
		RootDiskSizeGb: desired.EphemeralDiskSize,
	}, nil
}
