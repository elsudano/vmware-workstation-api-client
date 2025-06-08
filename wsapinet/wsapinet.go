package wsapinet

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
)

// New functon is just to create a new object NETClient to make the different calls at VmWare Workstation Pro
func New(httpcaller *httpclient.HTTPClient) NETService {
	return &NETManager{netclient: httpcaller}
}

func (netm *NETManager) LoadNICS(vm *wsapivm.MyVm) (*InfoNICS, error) {
	return nil, nil
}

func (netm *NETManager) CreateNIC(vm *wsapivm.MyVm, t string, vnet string) (*InfoNICS, error) {
	return nil, nil
}

func (netm *NETManager) UpdateNIC(vm *wsapivm.MyVm, inx int32, t string, vnet string) (*InfoNICS, error) {
	return nil, nil
}

func (netm *NETManager) DeleteNIC(vm *wsapivm.MyVm, inx int32) error {
	return nil
}
