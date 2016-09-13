package action

import (
	"fmt"

	"bosh-google-cpi/google/instance_service"
)

type DiskCloudProperties struct {
	DiskType string `json:"type,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

type Environment map[string]interface{}

type NetworkCloudProperties struct {
	NetworkName         string        `json:"network_name,omitempty"`
	SubnetworkName      string        `json:"subnetwork_name,omitempty"`
	Tags                instance.Tags `json:"tags,omitempty"`
	EphemeralExternalIP bool          `json:"ephemeral_external_ip,omitempty"`
	IPForwarding        bool          `json:"ip_forwarding,omitempty"`
}

type SnapshotMetadata struct {
	Deployment string `json:"deployment,omitempty"`
	Job        string `json:"job,omitempty"`
	Index      string `json:"index,omitempty"`
}

type StemcellCloudProperties struct {
	Name           string `json:"name,omitempty"`
	Version        string `json:"version,omitempty"`
	Infrastructure string `json:"infrastructure,omitempty"`
	SourceURL      string `json:"source_url,omitempty"`
}

type VMCloudProperties struct {
	Zone                string           `json:"zone,omitempty"`
	Name                string           `json:"name,omitempty"`
	MachineType         string           `json:"machine_type,omitempty"`
	CPU                 int              `json:"cpu,omitempty"`
	RAM                 int              `json:"ram,omitempty"`
	RootDiskSizeGb      int              `json:"root_disk_size_gb,omitempty"`
	RootDiskType        string           `json:"root_disk_type,omitempty"`
	AutomaticRestart    bool             `json:"automatic_restart,omitempty"`
	OnHostMaintenance   string           `json:"on_host_maintenance,omitempty"`
	Preemptible         bool             `json:"preemptible,omitempty"`
	ServiceAccount      VMServiceAccount `json:"service_account,omitempty"`
	ServiceScopes       VMServiceScopes  `json:"service_scopes,omitempty"`
	TargetPool          string           `json:"target_pool,omitempty"`
	BackendService      string           `json:"backend_service,omitempty"`
	Tags                instance.Tags    `json:"tags,omitempty"`
	EphemeralExternalIP *bool            `json:"ephemeral_external_ip,omitempty"`
	IPForwarding        *bool            `json:"ip_forwarding,omitempty"`
}

func (n VMCloudProperties) Validate() error {
	if err := n.Tags.Validate(); err != nil {
		return err
	}

	// A custom service account requires at least one scope
	if n.ServiceAccount != "" && len(n.ServiceScopes) == 0 {
		return fmt.Errorf("You must define at least one service scope if you define a service account. The scope 'https://www.googleapis.com/auth/cloud-platform' will allow your service account to access all of its available IAM permissions.")
	}

	return nil
}

type VMServiceScopes []string
type VMServiceAccount string
type VMMetadata map[string]string
