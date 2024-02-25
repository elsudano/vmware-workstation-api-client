package wsapiclient

import (
	"bytes"
	"encoding/json"
	"log"
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
	responseBody, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: GetAllVMs Message: I can't read the json structure %s", err)
		return nil, err
	}

	for vm, value := range vms {
		responseBody, err := c.httpRequest("vms/"+value.IdVM, "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			return nil, err
		}
		responseBody, err = c.httpRequest("vms/"+value.IdVM+"/power", "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&vms[vm])
		if err != nil {
			return nil, err
		}
		responseBody, err = c.httpRequest("vms/"+value.IdVM+"/params/displayName", "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&tmpparam)
		vms[vm].Denomination = tmpparam.Value
		if err != nil {
			return nil, err
		}
		responseBody, err = c.httpRequest("vms/"+value.IdVM+"/params/annotation", "GET", bytes.Buffer{})
		if err != nil {
			return nil, err
		}
		err = json.NewDecoder(responseBody).Decode(&tmpparam)
		vms[vm].Description = tmpparam.Value
		if err != nil {
			return nil, err
		}
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", vms)
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
	var vms []MyVm
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
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Encoding VM error %#v\n", err)
		return nil, err
	}
	// }}}
	// -------- Making the request in order to create the new vm --------- {{{
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("vms", "POST", *requestBody)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}
	// }}}
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't permit see this the another way
	// --------- Read the path and the ID of the vm in order to load in the function --------- {{{
	response, err = c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	for tempvm, value := range vms {
		if value.IdVM == vm.IdVM {
			vm = vms[tempvm]
		}
	}
	// }}}
	// --------- Preparing the next request in order to get the power status of the vm --------- {{{
	response, err = c.httpRequest("vms/"+vm.IdVM+"/power", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Request Error in power status %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body power status %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in power status %#v\n", err)
		return nil, err
	}
	// }}}
	// --------- We will change the values of the settings on the VM  --------- {{{
	requestBody.Reset()
	err = json.NewEncoder(requestBody).Encode(&tempSettingVM)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Settings error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, err = c.httpRequest("vms/"+vm.IdVM, "PUT", *requestBody)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Request Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	responseBody = new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
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
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Param error %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	// response, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change description %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	// ----- Now, we change the Description ----
	// tempDataParam.Name = "annotation"
	// tempDataParam.Value = d
	// requestBody.Reset()
	// err = json.NewEncoder(requestBody).Encode(&tempDataParam)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Encoding Param error %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	// response, err = c.httpRequest("vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:response raw %#v\n", response)
	// responseBody.Reset()
	// _, err = responseBody.ReadFrom(response)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	// log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Response Body in change description %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	// if err != nil {
	// 	log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: CreateVM Obj:Response Error in change description %#v\n", err)
	// 	return nil, err
	// }
	//}}}
	return &vm, err
}

// ReadVM method return the object MyVm with the ID indicate in i.
// Input: i: string with the ID of the VM, Return: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) ReadVM(i string) (*MyVm, error) {
	var vms []MyVm
	var vm MyVm
	var tmpparam ParamPayload
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't permit see this the another way
	// --------- Read the path and the ID of the vm in order to load in the function --------- {{{
	response, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
		}
	}
	// }}}
	// --------- Read the propierties of the VM in order to load --------- {{{
	response, err = c.httpRequest("vms/"+vm.IdVM, "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Request Error trying get information %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:response raw get information %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Response Error trying get information %#v\n", err)
		return nil, err
	}
	// }}
	// --------- Read the status of power of the vm --------- {{{
	response, err = c.httpRequest("vms/"+vm.IdVM+"/power", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Request Error in power status %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:Response Body power status %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Response Error in power status %#v\n", err)
		return nil, err
	}
	// }}
	// --------- The last part is read the denomination and description of the vm --------- {{{
	response, err = c.httpRequest("vms/"+vm.IdVM+"/params/displayName", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:trying get denomination %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:Response trying get denomination %#v\n", response)
	err = json.NewDecoder(response).Decode(&tmpparam)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Response Error trying get denomination %#v\n", err)
		return nil, err
	}
	vm.Denomination = tmpparam.Value
	response, err = c.httpRequest("vms/"+vm.IdVM+"/params/annotation", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:trying get description %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:Response trying get description %#v\n", response)
	err = json.NewDecoder(response).Decode(&tmpparam)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: ReadVM Obj:Response Error trying get description %#v\n", err)
		return nil, err
	}
	vm.Description = tmpparam.Value
	// }}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:VM %#v\n", vm)
	return &vm, nil
}

// UpdateVM method to update a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update, n: string with the denomination of VM
// d: string with the description of the VM, p: int with the number of processors
// m: int with the size of memory Output: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) UpdateVM(i string, n string, d string, p int, m int) (*MyVm, error) {
	var buffer bytes.Buffer
	// there prepare the body request
	request, err := json.Marshal(map[string]int{
		"processors": p,
		"memory":     m,
	})
	if err != nil {
		return nil, err
	}
	buffer.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: Request Body %#v\n", buffer.String())
	_, err = c.httpRequest("vms/"+i, "PUT", buffer)
	if err != nil {
		return nil, err
	}
	vm, err := GetVM(c, i)
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM before %#v\n", vm)
	// ---- here we have to implement the code to update de description and denomination
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM after %#v\n", vm)
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
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("vms/registration", "POST", *requestBody)
	if err != nil {
		return nil, err
	}
	// Piensa si tiene que ser en este punto en la parte de httpRequest
	// tienes que poner en este punto un control de errores de lo que responde VMW
	// si es diferente de create, ok, o delete que de un error y ponga a nil la vm y salga
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: RegisterVM Obj:Response Error %#v\n", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: RegisterVM Obj:Response Body %#v\n", responseBody.String())
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}
	return &vm, err
}

// DeleteVM method to delete a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update
func (c *Client) DeleteVM(i string) error {
	response, err := c.httpRequest("vms/"+i, "DELETE", bytes.Buffer{})
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", err)
		return err
	}
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapivm.go Fu: DeleteVM Obj:%#v, %#v\n", err, responseBody.String())
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: DeleteVM Obj:%#v\n", responseBody)
	return nil
}
