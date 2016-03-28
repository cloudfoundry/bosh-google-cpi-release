package action

type DiskCloudProperties struct {
	DiskType string `json:"type,omitempty"`
}

type Environment map[string]interface{}

type NetworkCloudProperties struct {
	NetworkName         string      `json:"network_name,omitempty"`
	SubnetworkName      string      `json:"subnetwork_name,omitempty"`
	Tags                NetworkTags `json:"tags,omitempty"`
	EphemeralExternalIP bool        `json:"ephemeral_external_ip,omitempty"`
	IPForwarding        bool        `json:"ip_forwarding,omitempty"`
	TargetPool          string      `json:"target_pool,omitempty"`
	InstanceGroup       string      `json:"instance_group,omitempty"`
}

type NetworkTags []string

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
	Zone              string          `json:"zone,omitempty"`
	MachineType       string          `json:"machine_type,omitempty"`
	CPU               int             `json:"cpu,omitempty"`
	RAM               int             `json:"ram,omitempty"`
	RootDiskSizeGb    int             `json:"root_disk_size_gb,omitempty"`
	RootDiskType      string          `json:"root_disk_type,omitempty"`
	AutomaticRestart  bool            `json:"automatic_restart,omitempty"`
	OnHostMaintenance string          `json:"on_host_maintenance,omitempty"`
	Preemptible       bool            `json:"preemptible,omitempty"`
	ServiceScopes     VMServiceScopes `json:"service_scopes,omitempty"`
}

type VMServiceScopes []string

type VMMetadata map[string]interface{}
