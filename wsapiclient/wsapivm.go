package wsapiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"strconv"
)

type MyVm struct {
	IdVM         string `json:"id"`
	Path         string `json:"path"`
	Denomination string `json:"displayName"`
	Description  string `json:"annotation"`
	// Image        string `json:"image"`
	CPU struct {
		Processors int `json:"processors"`
	}
	PowerStatus string `json:"power_state"`
	Memory      int    `json:"memory"`
	Ip          string `json:"ip"`
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

// GetAllVMs Method return array of MyVm and a error variable if occurr some problem
// Return: []MyVm and error
func (c *Client) GetAllVMs() ([]MyVm, error) {
	var vms []MyVm
	var tmpparam ParamPayload
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

	for vm, value := range vms {
		responseBody, vmerror, err := c.httpRequest("vms/"+value.IdVM, "GET", bytes.Buffer{})
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The request error: %#v", err)
			return nil, err
		}
		if vmerror.Code != 0 {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The 2 error API was %d %s", vmerror.Code, vmerror.Message)
			return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The Decoder error: %#v", err)
			return nil, err
		}
		responseBody, vmerror, err = c.httpRequest("vms/"+value.IdVM+"/power", "GET", bytes.Buffer{})
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The request error: %#v", err)
			return nil, err
		}
		if vmerror.Code != 0 {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The 3 error API was %d %s", vmerror.Code, vmerror.Message)
			return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The Decoder error: %#v", err)
			return nil, err
		}
		responseBody, vmerror, err = c.httpRequest("vms/"+value.IdVM+"/params/displayName", "GET", bytes.Buffer{})
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The request error: %#v", err)
			return nil, err
		}
		if vmerror.Code != 0 {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The 4 error API was %d %s", vmerror.Code, vmerror.Message)
			return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
		}
		err = json.NewDecoder(responseBody).Decode(&tmpparam)
		vms[vm].Denomination = tmpparam.Value
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The Decoder error: %#v", err)
			return nil, err
		}
		responseBody, vmerror, err = c.httpRequest("vms/"+value.IdVM+"/params/annotation", "GET", bytes.Buffer{})
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The request error: %#v", err)
			return nil, err
		}
		if vmerror.Code != 0 {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The 5 error API was %d %s", vmerror.Code, vmerror.Message)
			return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
		}
		err = json.NewDecoder(responseBody).Decode(&tmpparam)
		vms[vm].Description = tmpparam.Value
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs M: The Decoder error: %#v", err)
			return nil, err
		}
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
	// }}}
	// -------- Making the request in order to create the new vm --------- {{{
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
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
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
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
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
	// }}}
	// --------- This part read the Actual informations that we have about of the VM --------
	c.GetVM(vm.IdVM)
	// --------- We will change the values of the settings on the VM  --------- {{{
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
	if vmerror.Code != 0 {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM M: The 2 error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response RAW %#v\n", response)
	responseBody = new(bytes.Buffer)
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
	// }}}
	// The following code we will use in the future when the VmWare fix it the method configparams {{{
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
	// response, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change description %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
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
	// response, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change description %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	//}}}
	return &vm, err
}

// ReadVM method return the object MyVm with the ID indicate in i.
// Input: i: string with the ID of the VM, Return: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) ReadVM(i string) (*MyVm, error) {
	return c.GetVM(i)
}

// UpdateVM method to update a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update, n: string with the denomination of VM
// d: string with the description of the VM, p: int with the number of processors
// m: int with the size of memory, s: Power State desired
// Output: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) UpdateVM(i string, n string, d string, p int, m int, s string) (*MyVm, error) {
	var buffer bytes.Buffer
	// We want to know which is the current status of teh VM
	vm, err := c.GetVM(i)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Get Info Error %#v\n", err)
		return nil, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM before %#v\n", vm)
	/// Here We are preparing update the Power State of teh VM {{{
	if vm.PowerStatus != s {
		_, err = c.PowerSwitch(i, s)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Power Switch Error %#v\n", err)
			return nil, err
		}
	}
	// }}}
	currentPowerStatus := vm.PowerStatus
	// Here we are preparing the update of the Processors and Memory in the VM {{{
	if vm.CPU.Processors != p || vm.Memory != m {
		if currentPowerStatus != "off" {
			c.PowerSwitch(vm.IdVM, "off")
		}
		request, err := json.Marshal(map[string]int{
			"processors": p,
			"memory":     m,
		})
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
		if currentPowerStatus == "on" {
			c.PowerSwitch(vm.IdVM, "on")
		}
	}
	// }}}
	// ---- here we have to implement the code to update de description and denomination{{{
	// here you will need to use the API to change the values of the Denomination and Description
	// }}}
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

// GetNetwork Method to get all the Network information of the instance
// i: string with the ID of the VM to get Network information,
func (c *Client) GetNetwork(i string) (*MyVm, error) {
	var vm MyVm
	return &vm, nil
}

// PowerSwitch method that permit you change the state of the instance, so you will change
// from power-off to power-on the state of the instance.
// i: string with the ID of the VM to change the state,
// s: string with the state that will want between on, off, reset
func (c *Client) PowerSwitch(i string, s string) (*MyVm, error) {
	vm, err := c.GetVM(i)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch Obj: Error when Get VM %#v\n", err)
		return nil, err
	}
	requestBody := bytes.NewBufferString(s)
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch Obj: Request option %#v\n", requestBody.String())
	response, vmerror, err := c.httpRequest("vms/"+i+"/power", "PUT", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch Obj: Response RAW %#v\n", err)
		return nil, err
	}
	if vmerror.Code != 0 {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch Obj: Response Body RAW %#v, %#v\n", err, response)
		return nil, err
	}
	return vm, nil
}

// DeleteVM method to delete a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update
func (c *Client) DeleteVM(i string) error {
	c.PowerSwitch(i, "off")
	response, vmerror, err := c.httpRequest("vms/"+i, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", err)
		return err
	}
	if vmerror.Code != 0 {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM M: The 1 error API was %d %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v, %#v\n", err, responseBody.String())
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", responseBody.String())
	return nil
}
