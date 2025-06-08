package wsapinet

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/rs/zerolog/log"
)

// GetNetwork Method to get all the Network information of the instance
// Inputs:
// c: (Pointer) The client that we use to made the API calls.
// vmid: (string) That's the VM ID that we want to check the NICs information.
// Outputs:
// error (error) If we have some error we can handle it here.
func GetInfoNics(vmc *httpclient.HTTPClient, vmid string) error {
	var currentNIC InfoNICS
	response, err := vmc.ApiCall("vms/"+vmid+"/nicips", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	err = json.NewDecoder(response).Decode(&currentNIC)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	log.Debug().Msgf("VM: %#v", currentNIC)
	log.Info().Msg("We have read the Network Information.")
	return err
}

// RenewMAC Auxiliar function to renew the MAC address of the VM, as you know
// some operations can't be made by API, and for that reason we will need
// to delete, and recreate the NIC with the same parameters.
// Inputs:
// c: (Pointer) The client that we use to made the API calls.
// vmid: (string) The VM that we want to renew the MAC address.
// Outputs:
// error: (error) We can handle the errors here.
func RenewMAC(vmc *httpclient.HTTPClient, vmid string) error {
	var currentNIC InfoNICS
	var newNIC NicPayload
	requestBody := new(bytes.Buffer)
	// err := c.PowerSwitch(vm, "off")
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RenewMAC Obj: We can't shutdown the VM %#v\n", err)
	// 	return err
	// }
	response, err := vmc.ApiCall("vms/"+vmid+"/nic", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	err = json.NewDecoder(response).Decode(&currentNIC)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	_, err = vmc.ApiCall("vms/"+vmid+"/nic/"+fmt.Sprint(currentNIC.NICS[0].Index), "DELETE", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call to delete the NIC.")
		return err
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
		log.Error().Err(err).Msg("The NIC JSON is malformed.")
		return err
	}
	log.Debug().Msgf("Request RAW: %#v", requestBody.String())
	_, err = vmc.ApiCall("vms/"+vmid+"/nic", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call to create the NIC.")
		return err
	}
	log.Debug().Msgf("VM: %#v", currentNIC)
	log.Info().Msg("We have changed the MAC address.")
	return err
}
