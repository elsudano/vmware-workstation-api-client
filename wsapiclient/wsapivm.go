package wsapiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
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
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The request error: %#v", err)
		return nil, err
	}
	if vmerror.Code != 0 {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj: Response Body%#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: I can't read the json structure %s", err)
		return nil, err
	}
	for _, item := range vms {
		// --------- This Block read the ID of the VM --------- {{{
		vm, err := c.GetVM(item.IdVM)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAllVMs M: %#v\n", err)
			return nil, err
		}
		// }}}
		// --------- This Block read the propierties of the VM in order to load --------- {{{
		err = c.GetBasicInfo(vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAllVMs M: %#v\n", err)
			return nil, err
		}
		// }}}
		// --------- This Block read the status of power of the vm --------- {{{
		err = c.GetPowerStatus(vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAllVMs M: %#v\n", err)
			return nil, err
		}
		// }}}
		// --------- This block read the denomination and description of the vm --------- {{{
		err = c.GetDenominationDescription(vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAllVMs M: %#v\n", err)
			return nil, err
		}
		// }}}
		// --------- This Block read the IP information --------- {{{
		if vm.PowerStatus == "on" {
			err = c.GetNetwork(vm)
			if err != nil {
				log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAllVMs M: %#v\n", err)
				return nil, err
			}
		}
		// }}}
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj: List of VMs %#v\n", vms)
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
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Encoding VM error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, vmerror, err := c.httpRequest("vms", "POST", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
		return nil, err
	}
	switch vmerror.Code {
	// here we have to wait for unlock the SourceVM and then create the next one
	// keep in mind that maybe is better do that in the provider side
	case 147:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 107:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The SourceVM isn't powered off: %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 108:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The VM already exists: %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	case 109:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The SourceVM was locked: %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: 1 Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response RAW %#v\n", response)
	responseBody.Reset()
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj: Decode Error %#v\n", err)
		return nil, err
	}
	requestBody.Reset()
	err = json.NewEncoder(requestBody).Encode(&tempSettingVM)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Settings error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM, "PUT", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
		return nil, err
	}
	switch vmerror.Code {
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: 2 Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response RAW %#v\n", response)
	responseBody.Reset()
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj: Decoder Error %#v\n", err)
		return nil, err
	}
	err = c.RenewMAC(&vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj: RenewMAC Error %#v\n", err)
		return nil, err
	}
	// ----- Now, we change the Denomination ----
	// tempDataParam.Name = "displayName"
	// tempDataParam.Value = n
	// requestBody.Reset()
	// err = json.NewEncoder(requestBody).Encode(&tempDataParam)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Param error %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	// response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
	// 	return nil, err
	// }
	// switch vmerror.Code {
	// default:
	// 	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: 3 Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Error changing denomination %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change denomination %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change denomination %#v\n", err)
	// 	return nil, err
	// }
	// ----- Now, we change the Description ----
	// tempDataParam.Name = "annotation"
	// tempDataParam.Value = d
	// requestBody.Reset()
	// err = json.NewEncoder(requestBody).Encode(&tempDataParam)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Param error %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	// response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
	// 	return nil, err
	// }
	// switch vmerror.Code {
	// default:
	// 	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: 3 Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Error changing description %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change description %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
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
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: LoadVM M: GetVM %#v\n", err)
		return nil, err
	}
	err = c.GetBasicInfo(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: LoadVM M: GetBasicInfo %#v\n", err)
		return nil, err
	}
	err = c.GetDenominationDescription(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: LoadVM M: GetDenominationDescription %#v\n", err)
		return nil, err
	}
	err = c.GetPowerStatus(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: LoadVM M: GetPowerStatus %#v\n", err)
		return nil, err
	}
	if vm.PowerStatus == "on" {
		err = c.GetNetwork(vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: LoadVM M: PowerStatus %#v\n", err)
			return nil, err
		}
	}
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
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Get Info Error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM before %#v\n", vm)
	if s == "" {
		currentPowerStatus = vm.PowerStatus
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: The Curret Power Status was %#v\n", currentPowerStatus)
	} else {
		currentPowerStatus = s
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: We want to change the current Power Status at %#v\n", currentPowerStatus)
	}
	// Here we are preparing the update of the Processors and Memory in the VM {{{
	err = c.PowerSwitch(vm, "off")
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: We can't shutdown the VM %#v\n", err)
		return nil, err
	}
	request, err := json.Marshal(memcpu)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Error Marsherling Body %#v\n", err)
		return nil, err
	}
	buffer.Write(request)
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Request Body %#v\n", buffer.String())
	_, vmerror, err := c.httpRequest("vms/"+i, "PUT", buffer)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Request Error %#v\n", err)
		return nil, err
	}
	switch vmerror.Code {
	// here we have to wait for unlock the SourceVM and then create the next one
	// keep in mind that maybe is better do that in the provider side
	case 147:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	err = c.PowerSwitch(vm, currentPowerStatus)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: We can't reestablishment the power status %#v\n", err)
		return nil, err
	}
	// ---- here we have to implement the code to update de description and denomination {{{
	// here you will need to use the API to change the values of the Denomination and Description
	// }}}
	err = c.GetBasicInfo(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM M: GetBasicInfo %#v\n", err)
		return nil, err
	}
	err = c.GetDenominationDescription(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM M: GetDenominationDescription %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM after %#v\n", vm)
	return vm, err
}

// RegisterVM method to register a new VM in VmWare Worstation GUI:
// n: string with the VM NAME, p: string with the path of the VM
func (c *Client) RegisterVM(n string, p string) (*MyVm, error) {
	var vm MyVm
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(map[string]string{
		"name": n,
		"path": p,
	})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM M: Error preparing the request %#v\n", request)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj: Body Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj: Prepared Request Body %#v\n", requestBody.String())
	response, vmerror, err := c.httpRequest("vms/registration", "POST", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM M: Requested Body %#v\n", requestBody.String())
		return nil, err
	}
	if vmerror.Code != 0 {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj: Decode Error %#v\n", err)
		return nil, err
	}
	return &vm, err
}

// DeleteVM method to delete a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update
func (c *Client) DeleteVM(i string) error {
	vm, err := c.LoadVM(i)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj: Get %#v\n", err)
		return err
	}
	err = c.PowerSwitch(vm, "off")
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj: Power %#v\n", err)
		return err
	}
	response, vmerror, err := c.httpRequest("vms/"+i, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		responseBody := new(bytes.Buffer)
		_, err = responseBody.ReadFrom(response)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v, %#v\n", err, responseBody.String())
			return err
		}
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", responseBody.String())
	case 107:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM M: Shutdown the VM %d %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetNetwork Obj: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	return nil
}
