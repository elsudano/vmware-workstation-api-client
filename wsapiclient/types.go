package wsapiclient

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapinet"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
)

// Interface with all the methods that we can use to talk with the API of VmWare Workstation Pro
type WSAPIService interface {
	ConfigLog(lvl string, mode string)
	ConfigApiClient(a string, u string, p string, i bool, d string) error
	GetAllVMs() ([]wsapivm.MyVm, error)
	LoadVM(i string) (*wsapivm.MyVm, error)
	LoadVMbyName(n string) (*wsapivm.MyVm, error)
	CreateVM(pid string, n string, d string, p int32, m int32, s string) (*wsapivm.MyVm, error)
	UpdateVM(vm *wsapivm.MyVm, n string, d string, p int32, m int32, s string) error
	RegisterVM(vm *wsapivm.MyVm) error
	DeleteVM(vm *wsapivm.MyVm) error
}

// That's the abstract object that we will use to interact with API of VmWare Workstation Pro
type WSAPIClient struct {
	Caller     *httpclient.HTTPClient
	VMService  wsapivm.VMService
	NETService wsapinet.NETService
}
