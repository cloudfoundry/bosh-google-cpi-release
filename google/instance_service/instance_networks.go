package instance

type Networks map[string]Network

type Network struct {
	Type                string
	IP                  string
	Gateway             string
	Netmask             string
	DNS                 []string
	Default             []string
	NetworkName         string
	Tags                NetworkTags
	EphemeralExternalIP bool
	IPForwarding        bool
	TargetPool          string
}

type NetworkTags []string

func (n Network) IsDynamic() bool { return n.Type == "dynamic" }

func (n Network) IsVip() bool { return n.Type == "vip" }
