package action

import (
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
	NoCPUErr = "CPU must be greater than 0"
)

func (CalculateVMCloudProperties) Run(desired DesiredVMSpec) (VMCloudProperties, error) {
	if desired.CPU <= 0 {
		return VMCloudProperties{}, bosherr.Error(NoCPUErr)
	}

	minRam := int(float64(desired.CPU) * float64(minMemoryPerCPU))
	if desired.RAM < minRam {
		desired.RAM = minRam
	}
	remainder := desired.RAM % memoryGranularity
	if remainder != 0 {
		desired.RAM += memoryGranularity - remainder
	}

	return VMCloudProperties{
		CPU:            desired.CPU,
		RAM:            desired.RAM,
		RootDiskSizeGb: desired.EphemeralDiskSize,
	}, nil
}
