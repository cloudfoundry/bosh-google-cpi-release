package instance

import (
	"google.golang.org/api/compute/v1"
)

type Service interface {
	AddAccessConfig(id string, zone string, networkInterface string, accessConfig *compute.AccessConfig) error
	AttachDisk(id string, diskLink string) (*DiskAttachmentDetail, error)
	AttachedDisks(id string) (AttachedDisks, error)
	DiskDetail(vmID string, diskLink string) (*DiskAttachmentDetail, error)
	CleanUp(id string)
	Create(vmProps *Properties, networks Networks, registryEndpoint string) (string, error, *compute.AttachedDisk)
	Delete(id string) error
	DeleteAccessConfig(id string, zone string, networkInterface string, accessConfig string) error
	DetachDisk(id string, diskID string) error
	Find(id string, zone string) (*compute.Instance, bool, error)
	Reboot(id string) error
	SetMetadata(id string, vmMetadata Metadata) error
	SetTags(id string, zone string, instanceTags *compute.Tags) error
	UpdateNetworkConfiguration(id string, networks Networks) error
}

type AttachedDisks []string

type Metadata map[string]string

type Properties struct {
	Zone              string
	Name              string
	Stemcell          string
	MachineType       string
	RootDiskSizeGb    int
	RootDiskType      string
	AutomaticRestart  bool
	OnHostMaintenance string
	Preemptible       bool
	ServiceAccount    ServiceAccount
	ServiceScopes     ServiceScopes
	TargetPool        string
	BackendService    BackendService
	Tags              Tags
	Labels            Labels
	Accelerators      []Accelerator
	EphemeralDiskType string
}

type ServiceScopes []string
type ServiceAccount string

type BackendService struct {
	Name string
}
type Accelerator struct {
	AcceleratorType string
	Count           int64
}
type DiskAttachmentDetail struct {
	Name string
	Path string
}
