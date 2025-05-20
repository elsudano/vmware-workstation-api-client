package wsapiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

// GetNetwork Method to get all the Network information of the instance
// i: string with the ID of the VM to get Network information,
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
		log.Printf("[ERROR][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Install VM-Tools Drivers %d %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	return err
}
