package wsapinet

import "github.com/elsudano/vmware-workstation-api-client/httpclient"

// New functon is just to create a new object NETClient to make the different calls at VmWare Workstation Pro
func New(httpcaller *httpclient.HTTPClient) NETService {
	return &NETManager{netclient: httpcaller}
}

func (netm *NETManager) CreateNIC(t string, vnet string) (*InfoNICS, error) { return nil, nil }

func (netm *NETManager) LoadNICS(vnet string) (*InfoNICS, error) { return nil, nil }

func (netm *NETManager) UpdateNIC(nic *InfoNICS, t string, vnet string) error { return nil }

func (netm *NETManager) DeleteNIC(nic *InfoNICS) error { return nil }
