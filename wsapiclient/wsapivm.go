package wsapiclient

import (
	"bytes"
	"encoding/json"
	"log"
)

type MyVm struct {
	IdVM         string `json:"id"`
	Path         string `json:"path"`
	Denomination string `json:"denomination"`
	Description  string `json:"description"`
	// Image        string `json:"image"`
	CPU struct {
		Processors int `json:"processors"`
	}
	PowerStatus string `json:"power_state"`
	Memory      int    `json:"memory"`
}

// Method GetAllVMs return array of MyVm and a error variable if occurr some problem
// Return: []MyVm and error
func (c *Client) GetAllVMs() ([]MyVm, error) {
	var vms []MyVm
	responseBody, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Message: I can't read the json structure %s", err)
		return nil, err
	}

	for vm, value := range vms {
		data, err := GetNameDescription(vms[vm].Path)
		if err != nil {
			return nil, err
		}
		vms[vm].Denomination = data[0]
		vms[vm].Description = data[1]

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
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:%#v\n", vms)
	return vms, nil
}

// Function to create a new VM in VmWare Worstation Input:
// s: string with the ID of the origin VM, d: string with the denomination of the VM
func (c *Client) CreateVM(s string, d string) (*MyVm, error) {
	var vm MyVm
	var requestBody bytes.Buffer
	request, err := json.Marshal(map[string]string{
		"name":     d,
		"parentId": s,
	})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: CreateVM Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("vms", "POST", requestBody)
	if err != nil {
		return nil, err
	}
	// Piensa si tiene que ser en este punto en la parte de httpRequest
	// tienes que poner en este punto un control de errores de lo que responde VMW
	// si es diferente de create, ok, o delete que de un error y ponga a nil la vm y salga
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

	// Falta hacer un PUT para modificar los parametros de la instancia nueva. entre ellos el procesador la memoria y la network
	return &vm, err
}

// Method ReadVM return the object MyVm with the ID indicate in i.
// Input: i: string with the ID of the VM, Return: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) ReadVM(i string) (*MyVm, error) {
	vm, err := GetVM(c, i)
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:VM %#v\n", vm)
	return vm, nil
}

// Function to update a VM in VmWare Worstation Input:
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
	err = SetNameDescription(vm.Path, n, d)
	if err != nil {
		return nil, err
	}
	vm, err = GetVM(c, i)
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj: VM after %#v\n", vm)
	return vm, err
}

// Function to delete a VM in VmWare Worstation Input:
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
