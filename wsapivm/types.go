package wsapivm

type MyVm struct {
	// Image        string `json:"image"`
	IdVM         string `json:"id"`
	Path         string `json:"path"`
	Denomination string `json:"displayName"`
	Description  string `json:"annotation"`
	PowerStatus  string `json:"power_state"`
	Memory       int32  `json:"memory"`
	CPU          struct {
		Processors int32 `json:"processors"`
	}
	NICS []struct {
		Mac string   `json:"mac"`
		Ip  []string `json:"ip"`
	}
	DNS struct {
		Hostname   string   `json:"hostname"`
		Domainname string   `json:"domainname"`
		Servers    []string `json:"server"`
	}
}

// This struct is for create a VM, just for create because the API needs
type CreatePayload struct {
	Name     string `json:"name"`
	ParentId string `json:"parentId"`
}

// This struct is for get and put the definition of VM
type SettingPayload struct {
	Processors int `json:"processors"`
	Memory     int `json:"memory"`
}

// I we want to register the VM in the GUI we will use this payload
type RegisterPayload struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

// This struct is for get and put information about of any parameters of the VM
type ParamPayload struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// This struct is for get and put information about of any Power State of the VM
type PowerStatePayload struct {
	Value string `json:"power_state"`
}
