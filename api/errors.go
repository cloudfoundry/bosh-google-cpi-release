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

type NotSupportedError struct{}

func (e NotSupportedError) Type() string  { return "Bosh::Clouds::NotSupported" }
func (e NotSupportedError) Error() string { return "Not supported" }

type VMNotFoundError struct {
	vmID string
}

func NewVMNotFoundError(vmID string) VMNotFoundError {
	return VMNotFoundError{vmID: vmID}
}

func (e VMNotFoundError) Type() string  { return "Bosh::Clouds::VMNotFound" }
func (e VMNotFoundError) Error() string { return fmt.Sprintf("VM '%s' not found", e.vmID) }

type VMCreationFailedError struct{}

func (e VMCreationFailedError) Type() string   { return "Bosh::Clouds::VMCreationFailed" }
func (e VMCreationFailedError) Error() string  { return "VM failed to create" }
func (e VMCreationFailedError) CanRetry() bool { return false }

type NoDiskSpaceError struct{}

func (e NoDiskSpaceError) Type() string   { return "Bosh::Clouds::NoDiskSpace" }
func (e NoDiskSpaceError) Error() string  { return "No disk space" }
func (e NoDiskSpaceError) CanRetry() bool { return false }

type DiskNotAttachedError struct {
	vmID   string
	diskID string
}

func NewDiskNotAttachedError(vmID, diskID string) DiskNotAttachedError {
	return DiskNotAttachedError{vmID: vmID, diskID: diskID}
}

func (e DiskNotAttachedError) Type() string { return "Bosh::Clouds::DiskNotAttached" }

func (e DiskNotAttachedError) Error() string {
	return fmt.Sprintf("Disk '%s' not attached to VM '%s'", e.diskID, e.vmID)
}

func (e DiskNotAttachedError) CanRetry() bool { return false }

type DiskNotFoundError struct {
	diskID string
}

func NewDiskNotFoundError(diskID string) DiskNotFoundError {
	return DiskNotFoundError{diskID: diskID}
}

func (e DiskNotFoundError) Type() string   { return "Bosh::Clouds::DiskNotFound" }
func (e DiskNotFoundError) Error() string  { return fmt.Sprintf("Disk '%s' not found", e.diskID) }
func (e DiskNotFoundError) CanRetry() bool { return false }
