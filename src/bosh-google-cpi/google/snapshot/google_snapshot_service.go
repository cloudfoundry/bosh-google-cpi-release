package snapshot

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
	"google.golang.org/api/compute/v1"

	"bosh-google-cpi/google/operation"
)

const googleSnapshotServiceLogTag = "GoogleSnapshotService"
const googleSnapshotNamePrefix = "snapshot"
const googleSnapshotDescription = "Snapshot managed by BOSH"
const googleSnapshotReadyStatus = "READY"
const googleSnapshotFailedStatus = "FAILED"

type GoogleSnapshotService struct {
	project          string
	computeService   *compute.Service
	operationService operation.Service
	uuidGen          boshuuid.Generator
	logger           boshlog.Logger
}

func NewGoogleSnapshotService(
	project string,
	computeService *compute.Service,
	operationService operation.Service,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) GoogleSnapshotService {
	return GoogleSnapshotService{
		project:          project,
		computeService:   computeService,
		operationService: operationService,
		uuidGen:          uuidGen,
		logger:           logger,
	}
}
