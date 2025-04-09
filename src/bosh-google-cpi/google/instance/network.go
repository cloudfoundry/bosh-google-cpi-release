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
	NetworkProjectID    string
	SubnetworkName      string
	EphemeralExternalIP bool
	IPForwarding        bool
	Tags                Tags
}

type Tags []string

func (t Tags) Validate() error {
	if len(t) > 0 {
		pattern := regexp.MustCompile("^[A-Za-z]+[A-Za-z0-9-]*[A-Za-z0-9]+$")
		for _, tag := range t {
			if len(tag) > maxTagLength || !pattern.MatchString(tag) {
				return bosherr.Errorf("Invalid tag '%s': does not comply with RFC1035", tag)
			}
		}
	}

	return nil
}

func (t Tags) Unique() []string {
	tagDict := make(map[string]struct{})
	for _, tag := range t {
		tagDict[tag] = struct{}{}
	}

	tagItems := make([]string, 0)
	for tag := range tagDict {
		tagItems = append(tagItems, tag)
	}
	return tagItems
}

func (n Network) IsDynamic() bool { return n.Type == "dynamic" }

func (n Network) IsVip() bool { return n.Type == "vip" }

func (n Network) IsManual() bool { return n.Type == "" || n.Type == "manual" }

func (n Network) Validate() error {
	switch {
	case n.IsDynamic():
		if err := n.Tags.Validate(); err != nil {
			return err
		}
	case n.IsManual():
		if err := n.Tags.Validate(); err != nil {
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
