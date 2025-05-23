package wsapiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/rs/zerolog/log"
)

type MyVm struct {
	// Image        string `json:"image"`
	IdVM         string `json:"id"`
	Path         string `json:"path"`
	Denomination string `json:"displayName"`
	Description  string `json:"annotation"`
	PowerStatus  string `json:"power_state"`
	Memory       int    `json:"memory"`
	CPU          struct {
		Processors int `json:"processors"`
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

// GetAllVMs Method return array of MyVm and a error variable if occurr some problem
// Outputs:
// []MyVm list of all VMs that we have in VmWare Workstation
// (error) variable with the error if occurr
func (c *Client) GetAllVMs() ([]MyVm, error) {
	var vms []MyVm
	responseBody, vmerror, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We can't made the API call.")
		return nil, err
	}
	if vmerror.Code != 0 {
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New("Code: " + strconv.Itoa(vmerror.Code) + ", Message: " + vmerror.Message)
	}
	log.Debug().Msgf("Response Body RAW: %#v", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Error().Err(err).Msg("The JSON was malformed")
		return nil, err
	}
	log.Info().Str("NumOfVMs", strconv.Itoa(len(vms))).Msg("You have this amount of VM in you Workstation")
	for _, item := range vms {
		// --------- This Block read the ID of the VM --------- {{{
		vm, err := c.LoadVM(item.IdVM)
		if err != nil {
			log.Error().Err(err).Msg("We can't Load the VM.")
			return nil, err
		}
		// }}}
		// --------- This Block read the propierties of the VM in order to load --------- {{{
		err = c.GetBasicInfo(vm)
		if err != nil {
			log.Error().Err(err).Msg("We can't read Basic Information")
			return nil, err
		}
		// }}}
		// --------- This Block read the status of power of the vm --------- {{{
		err = c.GetPowerStatus(vm)
		if err != nil {
			log.Error().Err(err).Msg("We can't read Power Status.")
			return nil, err
		}
		// }}}
		// --------- This block read the denomination and description of the vm --------- {{{
		err = c.GetDenominationDescription(vm)
		if err != nil {
			log.Error().Err(err).Msg("We can't read Description.")
			return nil, err
		}
		// }}}
		// --------- This Block read the IP information --------- {{{
		if vm.PowerStatus == "on" {
			err = c.GetNetwork(vm)
			if err != nil {
				log.Error().Err(err).Msg("We can't read Network Information.")
				return nil, err
			}
		}
		// }}}
	}
	log.Info().Msg("We have listed all VMs")
	return vms, nil
}

// CreateVM method to create a new VM in VmWare Worstation Input:
// s: string with the ID of the origin VM,
// n: string with the denomination of the VM,
// d: string with the description of VM
// p: int with the number of processors in the VM
// m: int with the number of memory in the VM
func (c *Client) CreateVM(s string, n string, d string, p int, m int) (*MyVm, error) {
	// --------- Preparing the request --------- {{{
	var vm MyVm
	requestBody := new(bytes.Buffer)
	responseBody := new(bytes.Buffer)
	var tempDataVM CreatePayload
	tempDataVM.Name = n
	tempDataVM.ParentId = s
	var tempSettingVM SettingPayload
	tempSettingVM.Processors = p
	tempSettingVM.Memory = m
	// var tempDataParam ParamPayload
	err := json.NewEncoder(requestBody).Encode(&tempDataVM)
	log.Debug().Msgf("Request Body RAW: %#v", requestBody.String())
	if err != nil {
		log.Error().Err(err).Msg("The request JSON is malformed.")
		return nil, err
	}
	response, vmerror, err := c.httpRequest("vms", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We can't made the API call.")
		return nil, err
	}
	switch vmerror.Code {
	// here we have to wait for unlock the SourceVM and then create the next one
	// keep in mind that maybe is better do that in the provider side
	case 0:
	case 147:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 107:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 108:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 109:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Debug().Msgf("Response RAW: %#v", response)
	responseBody.Reset()
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return nil, err
	}
	log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return nil, err
	}
	requestBody.Reset()
	err = json.NewEncoder(requestBody).Encode(&tempSettingVM)
	if err != nil {
		log.Error().Err(err).Msg("The Settings JSON is malformed.")
		return nil, err
	}
	log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM, "PUT", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return nil, err
	}
	switch vmerror.Code {
	case 0:
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Debug().Msgf("Response RAW: %#v", response)
	responseBody.Reset()
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't read the data from the Request.")
		return nil, err
	}
	log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return nil, err
	}
	err = c.RenewMAC(&vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the MAC information.")
		return nil, err
	}
	// ----- Now, we change the Denomination ----
	// tempDataParam.Name = "displayName"
	// tempDataParam.Value = n
	// requestBody.Reset()
	// err = json.NewEncoder(requestBody).Encode(&tempDataParam)
	// if err != nil {
	//	log.Error().Err(err).Msg("The Denomination JSON is malformed.")
	// 	return nil, err
	// }
	// log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	// response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	//	log.Error().Err(err).Msg("We couldn't read the data from the Request.")
	// 	return nil, err
	// }
	// switch vmerror.Code {
	// default:
	//	log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	// }
	// log.Debug().Msgf("Response RAW: %#v", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Error().Err(err).Msg("We couldn't read the data from the Response.")
	// 	return nil, err
	// }
	// log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	//	log.Error().Err(err).Msg("The response JSON is malformed.")
	// 	return nil, err
	// }
	// ----- Now, we change the Description ----
	// tempDataParam.Name = "annotation"
	// tempDataParam.Value = d
	// requestBody.Reset()
	// err = json.NewEncoder(requestBody).Encode(&tempDataParam)
	// if err != nil {
	//	log.Error().Err(err).Msg("The Description JSON is malformed.")
	// 	return nil, err
	// }
	// log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	// response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	//	log.Error().Err(err).Msg("We couldn't read the data from the Request.")
	// 	return nil, err
	// }
	// switch vmerror.Code {
	// default:
	//	log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	// }
	// log.Debug().Msgf("Response RAW: %#v", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Error().Err(err).Msg("We couldn't read the data from the Response.")
	// 	return nil, err
	// }
	// log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	//	log.Error().Err(err).Msg("The response JSON is malformed.")
	// 	return nil, err
	// }
	log.Info().Msg("We have created the VM.")
	return &vm, err
}

// ReadVM method return the object MyVm with the ID indicate in i.
// Inputs:
// i: (string) String with the ID of the VM
// Outputs:
// (pointer) Pointer at the MyVm object
// (error) variable with the error if occurr
func (c *Client) LoadVM(i string) (*MyVm, error) {
	vm, err := c.GetVM(i)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the ID and Path.")
		return nil, err
	}
	err = c.GetBasicInfo(vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Process and Memory.")
		return nil, err
	}
	err = c.GetDenominationDescription(vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Denomination and Description.")
		return nil, err
	}
	err = c.GetPowerStatus(vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Power Status.")
		return nil, err
	}
	if vm.PowerStatus == "on" {
		err = c.GetNetwork(vm)
		if err != nil {
			log.Error().Err(err).Msg("We couldn't show the Network information.")
			return nil, err
		}
	}
	log.Debug().Msgf("The ID that we are trying to load is: %#v", i)
	log.Info().Msg("We have loaded the VM.")
	return vm, err
}

// UpdateVM method to update a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update, n: string with the denomination of VM
// d: string with the description of the VM, p: int with the number of processors
// m: int with the size of memory
// s: Power State desired, choose between on, off, reset, (nil no change)
// Output: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) UpdateVM(i string, n string, d string, p int, m int, s string) (*MyVm, error) {
	var buffer bytes.Buffer
	var memcpu SettingPayload
	var currentPowerStatus string
	memcpu.Processors = p
	memcpu.Memory = m
	vm, err := c.LoadVM(i)
	// We want to know which is the current status of teh VM
	if err != nil {
		log.Error().Err(err).Msgf("We couldn't read the state of VM")
		return nil, err
	}
	log.Debug().Msgf("State of VM before to update: %#v", vm)
	if s == "" {
		currentPowerStatus = vm.PowerStatus
		log.Debug().Msgf("The Current Power Status was %#v", currentPowerStatus)
	} else {
		currentPowerStatus = s
		log.Debug().Msgf("We want to change the current Power Status at %#v", currentPowerStatus)
	}
	// Here we are preparing the update of the Processors and Memory in the VM {{{
	err = c.PowerSwitch(vm, "off")
	if err != nil {
		log.Error().Err(err).Msgf("We can't shutdown the VM")
		return nil, err
	}
	request, err := json.Marshal(memcpu)
	if err != nil {
		log.Error().Err(err).Msgf("Trying to encode this request: %#v", memcpu)
		return nil, err
	}
	buffer.Write(request)
	log.Debug().Msgf("Request Buffer: %#v", buffer.String())
	_, vmerror, err := c.httpRequest("vms/"+i, "PUT", buffer)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return nil, err
	}
	switch vmerror.Code {
	// here we have to wait for unlock the SourceVM and then create the next one
	// keep in mind that maybe is better do that in the provider side
	case 0:
	case 147:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	}
	err = c.PowerSwitch(vm, currentPowerStatus)
	if err != nil {
		log.Error().Err(err).Msgf("We can't complete the Shutdown/PowerOn operation")
		return nil, err
	}
	// ---- here we have to implement the code to update de description and denomination {{{
	// here you will need to use the API to change the values of the Denomination and Description
	// }}}
	err = c.GetBasicInfo(vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get information.")
		return nil, err
	}
	err = c.GetDenominationDescription(vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get information.")
		return nil, err
	}
	log.Debug().Msgf("State of VM after to update: %#v", vm)
	log.Info().Msg("We have updated the VM.")
	return vm, err
}

// RegisterVM method to register a new VM in VmWare Worstation GUI:
// n: string with the VM NAME, p: string with the path of the VM
func (c *Client) RegisterVM(n string, p string) (*MyVm, error) {
	var vm MyVm
	var regvm RegisterPayload
	regvm.Name = n
	regvm.Path = p
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(regvm)
	if err != nil {
		log.Error().Err(err).Msgf("Trying to encode this request: %#v", regvm)
		return nil, err
	}
	requestBody.Write(request)
	log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	response, vmerror, err := c.httpRequest("vms/registration", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return nil, err
	}
	switch vmerror.Code {
	case 0:
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("Response: %#v", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return nil, err
	}
	log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return nil, err
	}
	log.Info().Msg("We have registered the VM in GUI.")
	return &vm, err
}

// DeleteVM method to delete a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update
func (c *Client) DeleteVM(i string) error {
	vm, err := c.LoadVM(i)
	if err != nil {
		log.Error().Err(err).Msgf("We couldn't read the state of VM")
		return err
	}
	err = c.PowerSwitch(vm, "off")
	if err != nil {
		log.Error().Err(err).Msgf("We can't shutdown the VM")
		return err
	}
	response, vmerror, err := c.httpRequest("vms/"+i, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		responseBody := new(bytes.Buffer)
		_, err = responseBody.ReadFrom(response)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
		log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	case 107:
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Info().Msg("We have deleted the VM.")
	return nil
}
