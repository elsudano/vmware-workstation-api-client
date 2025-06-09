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
// vmc: (*httpclient.HTTPClient) The client that we use to made the API calls.
// vmid: (string) That's the VM ID that we want to check the NICs information.
// Outputs:
// NICS: (*InfoNICS) The structure with all the information about of the NICs that the VM has
// err: (error) If we have some error we can handle it here.
func GetNics(netc *httpclient.HTTPClient, vmid string) (NICS *InfoNICS, err error) {
	response, err := netc.ApiCall("vms/"+vmid+"/nic", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return NICS, err
	}
	err = json.NewDecoder(response).Decode(&NICS)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return NICS, err
	}
	log.Debug().Msgf("These's are the NIC's: %#v", NICS)
	log.Info().Msg("We have read the Network Information.")
	return NICS, nil
}

// CreateNic Auxiliary method to verate a NIC within VM
// vmc: (*httpclient.HTTPClient) The client that we use to made the API calls.
// vmid: (string) That's the VM ID that we want to check the NICs information.
// Outputs:
// NICS: (*InfoNICS) The structure with all the information about of the NICs that the VM has
// err: (error) If we have some error we can handle it here.
func CreateNic(netc *httpclient.HTTPClient, vmid string, t string, vnet string) (NIC *InfoNICS, err error) {
	var newNIC NicPayload
	requestBody := new(bytes.Buffer)
	if t == "bridged" {
		newNIC.Type = t
		newNIC.Vmnet = ""
	} else {
		newNIC.Type = t
		newNIC.Vmnet = vnet
	}
	err = json.NewEncoder(requestBody).Encode(&newNIC)
	if err != nil {
		log.Error().Err(err).Msg("The NIC JSON is malformed.")
		return nil, err
	}
	log.Debug().Msgf("Request RAW: %#v", requestBody.String())
	response, err := netc.ApiCall("vms/"+vmid+"/nic", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call creating NIC.")
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&NIC)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return NIC, err
	}
	log.Debug().Msgf("Info of new NIC: %#v", NIC)
	log.Info().Msg("We have created the NIC.")
	return NIC, nil
}

// DeleteNIC Auxiliary function to delete a NIC of one VM
// Inputs:
// vmc: (*httpclient.HTTPClient) The client that we use to made the API calls.
// vmid: (string) That's the VM ID that we want to check the NICs information.
// idx (int32) The array index of the NIC's has the VM
// Outputs:
// NICS: (*InfoNICS) The structure with all the information about of the NICs that the VM has.
// err: (error) If we have some error we can handle it here.
func DeleteNic(netc *httpclient.HTTPClient, vmid string, idx int32) (err error) {
	_, err = netc.ApiCall("vms/"+vmid+"/nic/"+fmt.Sprint(idx), "DELETE", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	log.Debug().Msgf("We have deleted this NIC: %#v", fmt.Sprint(idx))
	log.Info().Msg("We have Deleted the NIC.")
	return err
}

// RenewMAC Auxiliary function to renew the MAC address of the VM, as you know
// some operations can't be made by API, and for that reason we will need
// to delete, and recreate the NIC with the same parameters.
// Inputs:
// vmc: (*httpclient.HTTPClient) The client that we use to made the API calls.
// vmid: (string) The VM that we want to renew the MAC address.
// Outputs:
// error: (error) We can handle the errors here.
func RenewMAC(netc *httpclient.HTTPClient, vmid string) (err error) {
	var currentNIC InfoNICS
	var newNIC NicPayload
	requestBody := new(bytes.Buffer)
	// err := c.PowerSwitch(vm, "off")
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RenewMAC Obj: We can't shutdown the VM %#v\n", err)
	// 	return err
	// }
	response, err := netc.ApiCall("vms/"+vmid+"/nic", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	err = json.NewDecoder(response).Decode(&currentNIC)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	_, err = netc.ApiCall("vms/"+vmid+"/nic/"+fmt.Sprint(currentNIC.NICS[0].Index), "DELETE", bytes.Buffer{})
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
	_, err = netc.ApiCall("vms/"+vmid+"/nic", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call to create the NIC.")
		return err
	}
	log.Debug().Msgf("VM: %#v", currentNIC)
	log.Info().Msg("We have changed the MAC address.")
	return err
}
