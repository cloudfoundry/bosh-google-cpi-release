package api

import (
	"fmt"
)

type CloudError interface {
	error

	Type() string
}

type RetryableError interface {
	error

	CanRetry() bool
}

// -
type NotSupportedError struct{}

func (e NotSupportedError) Type() string  { return "Bosh::Clouds::NotSupported" }
func (e NotSupportedError) Error() string { return "Not supported" }

// -
type vmNotFoundError struct {
	vmID string
}

func NewVMNotFoundError(vmID string) vmNotFoundError {
	return vmNotFoundError{vmID: vmID}
}

func (e vmNotFoundError) Type() string  { return "Bosh::Clouds::VMNotFound" }
func (e vmNotFoundError) Error() string { return fmt.Sprintf("VM '%s' not found", e.vmID) }

// -
type VMCreationFailedError struct{}

func (e VMCreationFailedError) Type() string   { return "Bosh::Clouds::VMCreationFailed" }
func (e VMCreationFailedError) Error() string  { return "VM failed to create" }
func (e VMCreationFailedError) CanRetry() bool { return false }

// -
type NoDiskSpaceError struct{}

func (e NoDiskSpaceError) Type() string   { return "Bosh::Clouds::NoDiskSpace" }
func (e NoDiskSpaceError) Error() string  { return "No disk space" }
func (e NoDiskSpaceError) CanRetry() bool { return false }

// -
type diskNotAttachedError struct {
	vmID   string
	diskID string
}

func NewDiskNotAttachedError(vmID, diskID string) diskNotAttachedError {
	return diskNotAttachedError{vmID: vmID, diskID: diskID}
}

func (e diskNotAttachedError) Type() string { return "Bosh::Clouds::DiskNotAttached" }

func (e diskNotAttachedError) Error() string {
	return fmt.Sprintf("Disk '%s' not attached to VM '%s'", e.diskID, e.vmID)
}

func (e diskNotAttachedError) CanRetry() bool { return false }

// -
type diskNotFoundError struct {
	diskID string
}

func NewDiskNotFoundError(diskID string) diskNotFoundError {
	return diskNotFoundError{diskID: diskID}
}

func (e diskNotFoundError) Type() string   { return "Bosh::Clouds::DiskNotFound" }
func (e diskNotFoundError) Error() string  { return fmt.Sprintf("Disk '%s' not found", e.diskID) }
func (e diskNotFoundError) CanRetry() bool { return false }
