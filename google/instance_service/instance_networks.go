package ginstance

type InstanceNetworks map[string]InstanceNetwork

type InstanceNetwork struct {
	Type                string
	IP                  string
	Gateway             string
	Netmask             string
	DNS                 []string
	Default             []string
	NetworkName         string
	Tags                InstanceNetworkTags
	EphemeralExternalIP bool
	IPForwarding        bool
	TargetPool          string
}

type InstanceNetworkTags []string

func (n InstanceNetwork) IsDynamic() bool { return n.Type == "dynamic" }

func (n InstanceNetwork) IsVip() bool { return n.Type == "vip" }
