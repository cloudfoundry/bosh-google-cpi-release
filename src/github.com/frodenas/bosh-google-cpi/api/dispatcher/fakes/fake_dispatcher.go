package fakes

type FakeDispatcher struct {
	DispatchReqBytes  []byte
	DispatchRespBytes []byte
}

func (d *FakeDispatcher) Dispatch(reqBytes []byte) []byte {
	d.DispatchReqBytes = reqBytes
	return d.DispatchRespBytes
}
