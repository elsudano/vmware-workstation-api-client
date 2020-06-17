package wsapiclient

import (
	"bytes"
	"encoding/json"
	"log"
)

type MyVm struct {
	IdVM        string `json:"id"`
	Path        string `json:"path"`
	Description string `json:"description"`
	CPU         struct {
		Processors int32 `json:"processors"`
	}
	PowerStatus string `json:"power_state"`
	Memory      int32  `json:"memory"`
}

func (c *Client) GetAllVMs() ([]MyVm, error) {
	respBody, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Ob: %#v\n", respBody)
	var vms []MyVm
	err = json.NewDecoder(respBody).Decode(&vms)
	if err != nil {
		return nil, err
	}
	for id, value := range vms {
		respBody, err := c.httpRequest("vms/"+value.IdVM, "GET", bytes.Buffer{})
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(respBody).Decode(&vms[id])
		if err != nil {
			return nil, err
		}
		respBody, err = c.httpRequest("vms/"+value.IdVM+"/power", "GET", bytes.Buffer{})
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(respBody).Decode(&vms[id])
		if err != nil {
			return nil, err
		}
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Ob: %#v\n", vms)
	return vms, nil
}

func (c *Client) GetVM(idVM string) (*MyVm, error) {
	var vm MyVm
	body, err := c.httpRequest("vms/"+idVM, "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetVM Ob: %#v\n", body)
	err = json.NewDecoder(body).Decode(&vm)
	if err != nil {
		return nil, err
	}
	body, err = c.httpRequest("vms/"+idVM+"/power", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetVM Ob: %#v\n", body)
	err = json.NewDecoder(body).Decode(&vm)
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetVM Ob: %#v\n", vm)
	return &vm, nil
}
