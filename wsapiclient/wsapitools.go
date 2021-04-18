package wsapiclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	vmx "github.com/johlandabee/govmx"
)

// GetVM Auxiliar function to get the data of the VM and don't repeat code
// Input: c: pointer at the client of the API server, i: string with the ID yo VM
func GetVM(c *Client, i string) (*MyVm, error) {
	var vms []MyVm
	var vm MyVm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't permit see this the another way
	response, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	err = json.NewDecoder(response).Decode(&vms)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj: List of VMs %#v\n", vms)
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
		}
	}

	data, err := GetNameDescription(vm.Path)
	if err != nil {
		return nil, err
	}
	vm.Denomination = data[0]
	vm.Description = data[1]

	response, err = c.httpRequest("vms/"+i, "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj:Body of VM %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	response, err = c.httpRequest("vms/"+i+"/power", "GET", bytes.Buffer{})
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj:Body of power %#v\n", response)
	err = json.NewDecoder(response).Decode(&vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
		return nil, err
	}
	return &vm, nil
}

// GetNameDescription With this function you can see the Denomination and Description of the VM
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: p: string with the complete path of the file
// Output: []string with 2 positions, first denomination and second description
func GetNameDescription(p string) ([]string, error) {
	output := make([]string, 2)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetNameDescription Message: Failed opening file %s, please make sure the config file exists", err)
		return nil, err
	}

	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetNameDescription Obj: Data File %#v\n", string(data))

	vm := new(vmx.VirtualMachine)
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Fatalf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: GetNameDescription Obj: %#v", err)
		return nil, err
	}
	output[0] = vm.DisplayName
	output[1] = vm.Annotation
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetNameDescription Obj: VM  %#v\n", vm)
	return output, nil
}

// SetNameDescription With this function you can setting the Denomination and Description of the VM.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: p: string with the complete path of the file,
// n: string with the denomination, d: string with the description err: variable with error if occur
func SetNameDescription(p string, n string, d string) error {
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: parameters %#v, %#v, %#v", p, n, d)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed opening file %s, please make sure the config file exists", err)
		return err
	}

	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: File object %#v\n", string(data))

	vm := new(vmx.VirtualMachine)
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: GetNameDescription Obj: %#v", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: GetNameDescription Obj: VM %#v\n", vm)

	vm.DisplayName = n
	vm.Annotation = d
	data, err = vmx.Marshal(vm)
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: Data File %#v\n", string(data))
	err = ioutil.WriteFile(p, data, 0644)
	if err != nil {
		log.Fatalf("[WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed writing in file %s, please make sure the config file exists", err)
		return err
	}
	// en este punto tambien tienes que cambiar el nombre del fihero cuando se cambia la denominacion
	return err
}

// SetParameter With this function you can set the value of the parameter.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: i: string with the id of the VM,
// p: string with the name or param to set, v: string with the value of param err: variable with error if occur
func (c *Client) SetParameter(i string, p string, v string) error {
	requestBody := new(bytes.Buffer)
	request, err := json.Marshal(map[string]string{
		"name":  p,
		"value": v,
	})
	if err != nil {
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request %#v\n", request)
	requestBody.Write(request)
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request Body %#v\n", requestBody.String())
	response, err := c.httpRequest("/vms/"+i+"/configparams", "PUT", *requestBody)
	if err != nil {
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:response raw %#v\n", response)
	responseBody := new(bytes.Buffer)
	_, err = responseBody.ReadFrom(response)
	if err != nil {
		log.Printf("[WSAPICLI][ERROR] Fi: wsapitools.go Fu: SetParameter Obj:Response Error %#v\n", err)
		return err
	}
	log.Printf("[WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Response Body %#v\n", responseBody.String())
	// err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return err
	}
	return err
}
