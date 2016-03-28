package instance

import (
	"regexp"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

const maxTagLength = 63

type Network struct {
	Type                string
	IP                  string
	Gateway             string
	Netmask             string
	DNS                 []string
	Default             []string
	NetworkName         string
	SubnetworkName      string
	EphemeralExternalIP bool
	IPForwarding        bool
	Tags                NetworkTags
	TargetPool          string
	InstanceGroup       string
}

type NetworkTags []string

func (n Network) IsDynamic() bool { return n.Type == "dynamic" }

func (n Network) IsVip() bool { return n.Type == "vip" }

func (n Network) Validate() error {
	switch {
	case n.IsDynamic():
		if err := n.validateTags(); err != nil {
			return err
		}
	case n.IsVip():
		if n.IP == "" {
			return bosherr.Error("VIP Networks must provide an IP Address")
		}
	default:
		return bosherr.Errorf("Network type '%s' not supported", n.Type)
	}

	return nil
}

func (n Network) validateTags() error {
	if len(n.Tags) > 0 {
		pattern, _ := regexp.Compile("^[A-Za-z]+[A-Za-z0-9-]*[A-Za-z0-9]+$")
		for _, tag := range n.Tags {
			if len(tag) > maxTagLength || !pattern.MatchString(tag) {
				return bosherr.Errorf("Invalid tag '%s': does not comply with RFC1035", tag)
			}
		}
	}

	return nil
}
