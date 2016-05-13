package api

import (
	"bytes"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type MultiLogger struct {
	boshlog.Logger
	LogBuff *bytes.Buffer
}
