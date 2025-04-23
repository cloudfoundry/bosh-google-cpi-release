package action

import (
	"bosh-google-cpi/constant"
)

type InfoResult struct {
	StemcellFormats []string `json:"stemcell_formats"`
	ApiVersion      int      `json:"api_version,omitempty"`
}

type Info struct{}

func NewInfo() Info { return Info{} }

func (Info) Run() (InfoResult, error) {
	return InfoResult{
		StemcellFormats: []string{
			"google-light",
			"google-rawdisk",
		},
		ApiVersion: constant.MaxSupportedAPIVersion,
	}, nil
}
