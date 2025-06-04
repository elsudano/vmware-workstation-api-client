package wsapinet

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
