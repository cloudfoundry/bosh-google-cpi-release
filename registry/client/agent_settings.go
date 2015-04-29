package registry

const defaultNtpServer = "169.254.169.254"
const defaultSystemDisk = "/dev/sda"

type AgentSettingsResponse struct {
	Settings string `json:"settings"`
	Status   string `json:"status"`
}

type AgentSettings struct {
	AgentID   string            `json:"agent_id"`
	Blobstore BlobstoreSettings `json:"blobstore"`
	Disks     DisksSettings     `json:"disks"`
	Env       EnvSettings       `json:"env"`
	Mbus      string            `json:"mbus"`
	Networks  NetworksSettings  `json:"networks"`
	Ntp       []string          `json:"ntp"`
	VM        VMSettings        `json:"vm"`
}

type BlobstoreSettings struct {
	Provider string                 `json:"provider"`
	Options  map[string]interface{} `json:"options"`
}

type DisksSettings struct {
	System     string             `json:"system"`
	Persistent PersistentSettings `json:"persistent"`
}

type PersistentSettings map[string]string

type EnvSettings map[string]interface{}

type NetworksSettings map[string]NetworkSettings

type NetworkSettings struct {
	Type            string                 `json:"type"`
	IP              string                 `json:"ip"`
	Gateway         string                 `json:"gateway"`
	Netmask         string                 `json:"netmask"`
	DNS             []string               `json:"dns"`
	Default         []string               `json:"default"`
	CloudProperties map[string]interface{} `json:"cloud_properties"`
}

type VMSettings struct {
	ID string `json:"id"`
}

func NewAgentSettingsForVM(agentID string, vmCID string, networksSettings NetworksSettings, env EnvSettings, agentOptions AgentOptions) AgentSettings {
	var ntp []string
	if len(agentOptions.Ntp) > 0 {
		ntp = agentOptions.Ntp
	} else {
		ntp = append(ntp, defaultNtpServer)
	}

	agentSettings := AgentSettings{
		AgentID: agentID,
		Disks: DisksSettings{
			System:     defaultSystemDisk,
			Persistent: PersistentSettings{}},
		Blobstore: BlobstoreSettings{
			Provider: agentOptions.Blobstore.Type,
			Options:  agentOptions.Blobstore.Options,
		},
		Env:      EnvSettings(env),
		Mbus:     agentOptions.Mbus,
		Networks: networksSettings,
		Ntp:      ntp,
		VM: VMSettings{
			ID: vmCID,
		},
	}

	return agentSettings
}

func (as AgentSettings) AttachPersistentDisk(diskID string, deviceName string) AgentSettings {
	persistenSettings := PersistentSettings{}

	if as.Disks.Persistent != nil {
		for k, v := range as.Disks.Persistent {
			persistenSettings[k] = v
		}
	}

	persistenSettings[diskID] = deviceName
	as.Disks.Persistent = persistenSettings

	return as
}

func (as AgentSettings) ConfigureNetworks(networksSettings NetworksSettings) AgentSettings {
	as.Networks = networksSettings

	return as
}

func (as AgentSettings) DetachPersistentDisk(diskID string) AgentSettings {
	persistenSettings := PersistentSettings{}

	if as.Disks.Persistent != nil {
		for k, v := range as.Disks.Persistent {
			persistenSettings[k] = v
		}
	}

	delete(persistenSettings, diskID)
	as.Disks.Persistent = persistenSettings

	return as
}
