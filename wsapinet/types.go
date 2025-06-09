package wsapinet

import (
	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
)

// Interface with all the methods that we can use to talk with the API of VmWare Workstation Pro
type NETService interface {
	CreateNIC(vm *wsapivm.MyVm, t string, vnet string) (*InfoNICS, error)
	LoadNICS(vm *wsapivm.MyVm) (*InfoNICS, error)
	UpdateNIC(vm *wsapivm.MyVm, inx int32, t string, vnet string) (*InfoNICS, error)
	DeleteNIC(vm *wsapivm.MyVm, inx int32) error
}

// That's the Manager to make the calls
type NETManager struct {
	netclient *httpclient.HTTPClient
}

// This struct is to create a new NIC the information is different in APIRest of VmWare Workstation PRO
type NewNIC struct {
	Index int32  `json:"index"`
	Type  string `json:"type"`
	Vmnet string `json:"vmnet"`
	Mac   string `json:"macAddress"`
}

// This struct is for get and put information about NIC of the VM
type InfoNICS struct {
	Num  int `json:"num"`
	NICS []struct {
		Index int32  `json:"index"`
		Type  string `json:"type"`
		Vmnet string `json:"vmnet"`
		Mac   string `json:"macAddress"`
	}
}

// This is the information that we need to use in order to create a new NIC
type NicPayload struct {
	Type  string `json:"type"`
	Vmnet string `json:"vmnet"`
}

// This's the complete information of Network on a VM
type InfoNetwork struct {
	Nics []struct {
		Mac string   `json:"mac"`
		Ip  []string `json:"ip"`
		Dns struct {
			Hostname   string   `json:"hostname"`
			Domainname string   `json:"domainname"`
			Server     []string `json:"server"`
			Search     []string `json:"search"`
		}
		Wins struct {
			Primary   string `json:"primary"`
			Secondary string `json:"secondary"`
		}
		Dhcp4 struct {
			Enabled  bool   `json:"enabled"`
			Settings string `json:"settings"`
		}
		Dhcp6 struct {
			Enabled  bool   `json:"enabled"`
			Settings string `json:"settings"`
		}
	}
	Routes []struct {
		Dest      string `json:"dest"`
		Prefix    int32  `json:"prefix"`
		Nexthop   string `json:"nexthop"`
		Interface int32  `json:"interface"`
		Type      int32  `json:"type"`
		Metric    int32  `json:"metric"`
	}
	Dns struct {
		Hostname   string   `json:"hostname"`
		Domainname string   `json:"domainname"`
		Server     []string `json:"server"`
		Search     []string `json:"search"`
	}
	Wins struct {
		Primary   string `json:"primary"`
		Secondary string `json:"secondary"`
	}
	Dhcpv4 struct {
		Enabled  bool   `json:"enabled"`
		Settings string `json:"settings"`
	}
	Dhcpv6 struct {
		Enabled  bool   `json:"enabled"`
		Settings string `json:"settings"`
	}
}
