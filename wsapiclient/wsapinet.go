package wsapiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

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

// GetNetwork Method to get all the Network information of the instance
// Inputs:
// c: (Pointer) The client that we use to made the API calls.
// vm: (MyVM) That's the VM that we want to check the Network.
// Outputs:
// error (error) If we have some error we can handle it here.
func (c *Client) GetNetwork(vm *MyVm) error {
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/nicips", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj:  %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Requesting network information %#v\n", response)
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Error decoding Info Network %#v\n", err)
			return err
		}
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Response VM %#v\n", vm)
	case 118:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Output Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	}
	return err
}

// RenewMAC Auxiliar function to renew the MAC address of the VM, as you know
// some operations can't be made by API, and for that reason we will need
// to delete, and recreate the NIC with the same parameters.
// Inputs:
// c: (Pointer) The client that we use to made the API calls.
// vm: (MyVm) The VM that we want to change.
// Outputs:
// error: (error) We can handle the errors here.
func (c *Client) RenewMAC(vm *MyVm) error {
	var currentNIC InfoNICS
	var newNIC nicPayload
	requestBody := new(bytes.Buffer)
	// err := c.PowerSwitch(vm, "off")
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RenewMAC Obj: We can't shutdown the VM %#v\n", err)
	// 	return err
	// }
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/nic", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj:  %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj: Getting NIC %#v\n", response)
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&currentNIC)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj: Error decoding Info NIC %#v\n", err)
			return err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Output Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	}
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/nic/"+strconv.Itoa(currentNIC.NICS[0].Index), "DELETE", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj:  %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj: Deleting NIC %#v\n", response)
	switch vmerror.Code {
	case 0:
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Output Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	}
	if currentNIC.NICS[0].Type == "bridged" {
		newNIC.Type = currentNIC.NICS[0].Type
		newNIC.Vmnet = ""
	} else {
		newNIC.Type = currentNIC.NICS[0].Type
		newNIC.Vmnet = currentNIC.NICS[0].Vmnet
	}
	err = json.NewEncoder(requestBody).Encode(&newNIC)
	if err != nil {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj:%#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj:Request Body %#v\n", requestBody.String())
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/nic", "POST", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj:  %#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj: Creating NIC %#v\n", response)
	switch vmerror.Code {
	case 0:
	case 121:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Output Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: RenewMAC Obj: VM %#v\n", vm)
	return err
}
