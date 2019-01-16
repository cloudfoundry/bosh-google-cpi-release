package action

type InfoResult struct {
	StemcellFormats []string `json:"stemcell_formats"`
	ApiVersion int `json:"api_version,omitempty"`
}

type Info struct{}

func NewInfo() Info { return Info{} }

func (Info) Run() (InfoResult, error) {
	return InfoResult{
		StemcellFormats: []string{
			"google-light",
			"google-rawdisk",
		},
		ApiVersion: 1,
	}, nil
}
