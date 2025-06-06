package wsapinet

import "github.com/elsudano/vmware-workstation-api-client/httpclient"

type NETService interface {
	CreateNIC(t string, vnet string) (*InfoNICS, error)
	LoadNICS(vnet string) (*InfoNICS, error)
	UpdateNIC(nic *InfoNICS, t string, vnet string) error
	DeleteNIC(nic *InfoNICS) error
}

type NETManager struct {
	netclient *httpclient.HTTPClient
}

func New(httpcaller *httpclient.HTTPClient) NETService {
	return &NETManager{netclient: httpcaller}
}

func (netm *NETManager) CreateNIC(t string, vnet string) (*InfoNICS, error) { return nil, nil }

func (netm *NETManager) LoadNICS(vnet string) (*InfoNICS, error) { return nil, nil }

func (netm *NETManager) UpdateNIC(nic *InfoNICS, t string, vnet string) error { return nil }

func (netm *NETManager) DeleteNIC(nic *InfoNICS) error { return nil }
