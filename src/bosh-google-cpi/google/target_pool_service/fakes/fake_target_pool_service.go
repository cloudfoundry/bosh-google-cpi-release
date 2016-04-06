package fakes

import (
	"bosh-google-cpi/google/target_pool_service"
)

type FakeTargetPoolService struct {
	AddInstanceCalled bool
	AddInstanceErr    error

	FindCalled     bool
	FindFound      bool
	FindTargetPool targetpool.TargetPool
	FindErr        error

	FindByInstanceCalled     bool
	FindByInstanceFound      bool
	FindByInstanceTargetPool string
	FindByInstanceErr        error

	ListCalled      bool
	ListTargetPools []targetpool.TargetPool
	ListErr         error

	RemoveInstanceCalled bool
	RemoveInstanceErr    error
}

func (t *FakeTargetPoolService) AddInstance(id string, vmLink string) error {
	t.AddInstanceCalled = true
	return t.AddInstanceErr
}

func (t *FakeTargetPoolService) Find(id string, region string) (targetpool.TargetPool, bool, error) {
	t.FindCalled = true
	return t.FindTargetPool, t.FindFound, t.FindErr
}

func (t *FakeTargetPoolService) FindByInstance(vmLink string, region string) (string, bool, error) {
	t.FindByInstanceCalled = true
	return t.FindByInstanceTargetPool, t.FindByInstanceFound, t.FindByInstanceErr
}

func (t *FakeTargetPoolService) List(region string) ([]targetpool.TargetPool, error) {
	t.ListCalled = true
	return t.ListTargetPools, t.ListErr
}

func (t *FakeTargetPoolService) RemoveInstance(id string, vmLink string) error {
	t.RemoveInstanceCalled = true
	return t.RemoveInstanceErr
}
