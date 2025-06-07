package wsapinet

import "github.com/elsudano/vmware-workstation-api-client/httpclient"

// Interface with all the methods that we can use to talk with the API of VmWare Workstation Pro
type NETService interface {
	CreateNIC(t string, vnet string) (*InfoNICS, error)
	LoadNICS(vnet string) (*InfoNICS, error)
	UpdateNIC(nic *InfoNICS, t string, vnet string) error
	DeleteNIC(nic *InfoNICS) error
}

// That's the Manager to make the calls
type NETManager struct {
	netclient *httpclient.HTTPClient
}

// This struct is for get and put information about NIC of the VM
type InfoNICS struct {
	Num  int `json:"num"`
	NICS []struct {
		Index int    `json:"index"`
		Type  string `json:"type"`
		Vmnet string `json:"vmnet"`
		Mac   string `json:"macAddress"`
	}
}

// This is the information that we need to use in order to create a new NIC
type nicPayload struct {
	Type  string `json:"type"`
	Vmnet string `json:"vmnet"`
}
