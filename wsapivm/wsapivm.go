package wsapivm

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/rs/zerolog/log"
)

type VMService interface {
	GetAllVMs() ([]MyVm, error)
	LoadVM(i string) (*MyVm, error)
	LoadVMbyName(n string) (*MyVm, error)
	CreateVM(pid string, n string, d string, p int, m int) (*MyVm, error)
	UpdateVM(vm *MyVm, n string, d string, p int, m int, s string) error
	RegisterVM(vm *MyVm) error
	DeleteVM(vm *MyVm) error
}

type VMManager struct {
	vmclient *httpclient.HTTPClient
}

func New(httpcaller *httpclient.HTTPClient) VMService {
	return &VMManager{vmclient: httpcaller}
}

// GetAllVMs Method return array of MyVm and a error variable if occurr some problem
// Outputs:
// []MyVm list of all VMs that we have in VmWare Workstation
// (error) variable with the error if occurr
func (vmm *VMManager) GetAllVMs() ([]MyVm, error) {
	var vms []MyVm
	responseBody, err := vmm.vmclient.ApiCall("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We can't made the API call.")
		return nil, err
	}
	log.Debug().Msgf("Response Body RAW: %#v", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Error().Err(err).Msg("The JSON was malformed")
		return nil, err
	}
	log.Info().Str("NumOfVMs", strconv.Itoa(len(vms))).Msg("You have this amount of VM in you Workstation")
	for pos, item := range vms {
		// --------- This Block read the ID of the VM --------- {{{
		err = GetAllExtraParameters(vmm.vmclient, &item)
		if err != nil {
			log.Error().Err(err).Msg("We couldn't get all the extra parameters.")
			return nil, err
		}
		vms[pos] = item
		log.Debug().Msgf("The VM loaded is:: %#v", item)
	}
	log.Info().Msg("We have listed all VMs")
	return vms, nil
}

// CreateVM method to create a new VM in VmWare Worstation
// Input:
// pid: (string) with the ID of the Parent VM,
// n: string with the denomination of the VM,
// d: string with the description of VM
// p: int with the number of processors in the VM
// m: int with the number of memory in the VM
func (vmm *VMManager) CreateVM(pid string, n string, d string, p int, m int) (*MyVm, error) {
	// --------- Preparing the request --------- {{{
	var vm MyVm
	requestBody := new(bytes.Buffer)
	responseBody := new(bytes.Buffer)
	var tempDataVM CreatePayload
	tempDataVM.Name = n
	tempDataVM.ParentId = pid
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
	response, err := vmm.vmclient.ApiCall("vms", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We can't made the API call.")
		return nil, err
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
	response, err = vmm.vmclient.ApiCall("vms/"+vm.IdVM, "PUT", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return nil, err
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
	// err = wsapinet.RenewMAC(c, vm.IdVM)
	// if err != nil {
	// 	log.Error().Err(err).Msg("We couldn't show the MAC information.")
	// 	return nil, err
	// }
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

// LoadVM method return the object MyVm with the ID indicate in i.
// Inputs:
// i: (string) String with the ID of the VM
// Outputs:
// (pointer) Pointer at the MyVm object
// (error) variable with the error if occurr
func (vmm *VMManager) LoadVM(i string) (*MyVm, error) {
	vm, err := GetVM(vmm.vmclient, i)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the ID and Path.")
		return nil, err
	}
	err = GetAllExtraParameters(vmm.vmclient, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get all the extra parameters.")
		return nil, err
	}
	log.Debug().Msgf("The ID that we are trying to load is: %#v", i)
	log.Info().Msg("We have loaded the VM.")
	return vm, err
}

// LoadVMbyName method return the object MyVm with the Name indicate in n.
// Inputs:
// n: (string) String with the Name of the VM
// Outputs:
// (pointer) Pointer at the MyVm object
// (error) variable with the error if occurr
func (vmm *VMManager) LoadVMbyName(n string) (*MyVm, error) {
	vm, err := GetVMbyName(vmm.vmclient, n)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the ID and Path.")
		return nil, err
	}
	err = GetAllExtraParameters(vmm.vmclient, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get all the extra parameters.")
		return nil, err
	}
	log.Debug().Msgf("The ID that we are trying to load is: %#v", n)
	log.Info().Msg("We have loaded the VM.")
	return vm, err
}

// UpdateVM method to update a VM in VmWare Worstation
// Input:
// vm (*MyVm) The VM that we want to update
// n: string with the denomination of VM
// d: string with the description of the VM, p: int with the number of processors
// m: int with the size of memory
// s: Power State desired, choose between on, off, reset, (nil no change)
// Output:
// pointer at the MyVm object
// and error variable with the error if occurr
func (vmm *VMManager) UpdateVM(vm *MyVm, n string, d string, p int, m int, s string) error {
	var buffer bytes.Buffer
	var memcpu SettingPayload
	var currentPowerStatus string
	memcpu.Processors = p
	memcpu.Memory = m
	log.Debug().Msgf("State of VM before to update: %#v", vm)
	if s == "" {
		currentPowerStatus = vm.PowerStatus
		log.Debug().Msgf("The Current Power Status was %#v", currentPowerStatus)
	} else {
		currentPowerStatus = s
		log.Debug().Msgf("We want to change the current Power Status at %#v", currentPowerStatus)
	}
	// Here we are preparing the update of the Processors and Memory in the VM {{{
	err := PowerSwitch(vmm.vmclient, vm, "off")
	if err != nil {
		log.Error().Err(err).Msgf("We can't shutdown the VM")
		return err
	}
	request, err := json.Marshal(memcpu)
	if err != nil {
		log.Error().Err(err).Msgf("Trying to encode this request: %#v", memcpu)
		return err
	}
	buffer.Write(request)
	log.Debug().Msgf("Request Buffer: %#v", buffer.String())
	_, err = vmm.vmclient.ApiCall("vms/"+vm.IdVM, "PUT", buffer)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	err = PowerSwitch(vmm.vmclient, vm, currentPowerStatus)
	if err != nil {
		log.Error().Err(err).Msgf("We can't complete the Shutdown/PowerOn operation")
		return err
	}
	// ---- here we have to implement the code to update de description and denomination {{{
	// here you will need to use the API to change the values of the Denomination and Description
	// }}}
	err = GetBasicInfo(vmm.vmclient, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get information.")
		return err
	}
	err = GetDenominationDescription(vmm.vmclient, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't get information.")
		return err
	}
	log.Debug().Msgf("State of VM after to update: %#v", vm)
	log.Info().Msg("We have updated the VM.")
	return err
}

// RegisterVM method to register a new VM in VmWare Worstation GUI:
// Input:
// c: (*wsapiclient.Client) The client to make the call.
// vm: (*wsapivm.MyVM) The VM object that we want to delete.
// Output:
// error: (error) The possible error that you will have.
func (vmm *VMManager) RegisterVM(vm *MyVm) error {
	var regvm RegisterPayload
	regvm.Name = vm.Denomination
	regvm.Path = vm.Path
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(regvm)
	if err != nil {
		log.Error().Err(err).Msgf("Trying to encode this request: %#v", regvm)
		return err
	}
	requestBody.Write(request)
	log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	response, err := vmm.vmclient.ApiCall("vms/registration", "POST", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	log.Debug().Msgf("Response: %#v", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	log.Info().Msg("We have registered the VM in GUI.")
	return err
}

// DeleteVM method to delete a VM in VmWare Worstation
// Input:
// c: (*wsapiclient.Client) The client to make the call.
// vm: (*wsapivm.MyVM) The VM object that we want to delete.
// Output:
// error: (error) The possible error that you will have.
func (vmm *VMManager) DeleteVM(vm *MyVm) error {
	err := PowerSwitch(vmm.vmclient, vm, "off")
	if err != nil {
		log.Error().Err(err).Msgf("We can't shutdown the VM")
		return err
	}
	response, err := vmm.vmclient.ApiCall("vms/"+vm.IdVM, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Error().Err(err).Msg("The response JSON is malformed.")
		return err
	}
	log.Debug().Msgf("Response Human Readable: %#v", responseBody.String())
	log.Info().Msg("We have deleted the VM.")
	return nil
}
