package action

import (
	"encoding/json"
	"fmt"

	"bosh-google-cpi/google/instance"
)

type DiskCloudProperties struct {
	DiskType string `json:"type,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

type Environment map[string]interface{}

type NetworkCloudProperties struct {
	NetworkName         string        `json:"network_name,omitempty"`
	NetworkProjectID    string        `json:"xpn_host_project_id,omitempty"`
	SubnetworkName      string        `json:"subnetwork_name,omitempty"`
	Tags                instance.Tags `json:"tags,omitempty"`
	EphemeralExternalIP bool          `json:"ephemeral_external_ip,omitempty"`
	IPForwarding        bool          `json:"ip_forwarding,omitempty"`
}

type SnapshotMetadata struct {
	Deployment string `json:"deployment,omitempty"`
	Job        string `json:"job,omitempty"`
	Index      Number `json:"index,omitempty"`
}

type StemcellCloudProperties struct {
	Name           string `json:"name,omitempty"`
	Version        string `json:"version,omitempty"`
	Infrastructure string `json:"infrastructure,omitempty"`
	SourceURL      string `json:"source_url,omitempty"`

	// URL of an existing image (Image.SelfLink)
	ImageURL   string `json:"image_url,omitempty"`
	SourceSha1 string `json:"raw_disk_sha1,omitempty"`
}

type VMCloudProperties struct {
	Zone                string           `json:"zone,omitempty"`
	Name                string           `json:"name,omitempty"`
	MachineType         string           `json:"machine_type,omitempty"`
	CPU                 int              `json:"cpu,omitempty"`
	RAM                 int              `json:"ram,omitempty"`
	MachineSeries       string           `json:"machine_series,omitempty"`
	RootDiskSizeGb      int              `json:"root_disk_size_gb,omitempty"`
	RootDiskType        string           `json:"root_disk_type,omitempty"`
	AutomaticRestart    bool             `json:"automatic_restart,omitempty"`
	OnHostMaintenance   string           `json:"on_host_maintenance,omitempty"`
	Preemptible         bool             `json:"preemptible,omitempty"`
	ProvisioningModel   string           `json:"provisioning_model,omitempty"`
	ServiceAccount      VMServiceAccount `json:"service_account,omitempty"`
	ServiceScopes       VMServiceScopes  `json:"service_scopes,omitempty"`
	TargetPool          string           `json:"target_pool,omitempty"`
	BackendService      interface{}      `json:"backend_service,omitempty"`
	Tags                instance.Tags    `json:"tags,omitempty"`
	Labels              instance.Labels  `json:"labels,omitempty"`
	EphemeralExternalIP *bool            `json:"ephemeral_external_ip,omitempty"`
	EphemeralDiskType   string           `json:"ephemeral_disk_type,omitempty"`
	IPForwarding        *bool            `json:"ip_forwarding,omitempty"`
	Accelerators        []Accelerator    `json:"accelerators,omitempty"`
	NodeGroup           string           `json:"node_group,omitempty"`
}

func (n VMCloudProperties) Validate() error {
	if err := n.Tags.Validate(); err != nil {
		return err
	}

	if err := n.Labels.Validate(); err != nil {
		return err
	}

	if n.ProvisioningModel != "" {
		validModels := map[string]bool{
			"SPOT":     true,
			"STANDARD": true,
		}
		if !validModels[n.ProvisioningModel] {
			return fmt.Errorf("invalid provisioning_model '%s': must be 'SPOT' or 'STANDARD'", n.ProvisioningModel)
		}
	}

	return nil
}

type VMServiceScopes []string
type VMServiceAccount string
type VMMetadata map[string]string
type Accelerator struct {
	AcceleratorType string `json:"type,omitempty"`
	Count           int64  `json:"count,omitempty"`
}

type Number string

func (n *Number) UnmarshalJSON(b []byte) error {
	if b[0] != '"' {
		b = []byte(fmt.Sprintf("\"%s\"", b))
	}
	return json.Unmarshal(b, (*string)(n))
}
