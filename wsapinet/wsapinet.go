package wsapinet

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
)

// New functon is just to create a new object NETClient to make the different calls at VmWare Workstation Pro
func New(httpcaller *httpclient.HTTPClient) NETService {
	return &NETManager{netclient: httpcaller}
}

func (netm *NETManager) LoadNICS(vm *wsapivm.MyVm) (NICS *InfoNICS, err error) {
	return GetNics(netm.netclient, vm.IdVM)
}

func (netm *NETManager) CreateNIC(vm *wsapivm.MyVm, t string, vnet string) (NIC *InfoNICS, err error) {
	return CreateNic(netm.netclient, vm.IdVM, t, vnet)
}

func (netm *NETManager) UpdateNIC(vm *wsapivm.MyVm, idx int32, t string, vnet string) (NIC *InfoNICS, err error) {
	return nil, nil
}

func (netm *NETManager) DeleteNIC(vm *wsapivm.MyVm, idx int32) (err error) {
	return DeleteNic(netm.netclient, vm.IdVM, idx)
}
