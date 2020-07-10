package wsapiclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
)

type MyVm struct {
	IdVM         string `json:"id"`
	Path         string `json:"path"`
	Denomination string `json:"denomination"`
	Description  string `json:"description"`
	Image        string `json:"image"`
	CPU          struct {
		Processors int32 `json:"processors"`
	}
	PowerStatus string `json:"power_state"`
	Memory      int32  `json:"memory"`
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
		return nil, err
	}

	for id, value := range vms {
		file, err := os.Open(vms[id].Path)
		if err != nil {
			log.Fatalf("Failed opening file %s, please make sure the config file exists", err)
		}

		log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: GetAllVMs Obj:File object %#v\n", file)

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		var vardenomination, vardescription string
		for scanner.Scan() {
			temp := strings.SplitN(scanner.Text(), " = ", 2)
			vtemp := strings.ToLower(temp[0])
			if vtemp == "displayname" {
				vardenomination = strings.TrimSpace(temp[1])
				vardenomination = strings.TrimSuffix(vardenomination, "\"")
				vardenomination = strings.TrimPrefix(vardenomination, "\"")
			}
			if vtemp == "annotation" {
				vardescription = strings.TrimSpace(temp[1])
				vardescription = strings.TrimSuffix(vardescription, "\"")
				vardescription = strings.TrimPrefix(vardescription, "\"")
			}
		}
		vms[id].Denomination = vardenomination
		vms[id].Description = vardescription
		file.Close()

		responseBody, err := c.httpRequest("vms/"+value.IdVM, "GET", bytes.Buffer{})
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(responseBody).Decode(&vms[id])
		if err != nil {
			return nil, err
		}
		responseBody, err = c.httpRequest("vms/"+value.IdVM+"/power", "GET", bytes.Buffer{})
		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(responseBody).Decode(&vms[id])
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
	// tienes que poner en este punto un control de errores de lo que responde VMW
	// si es diferente de create o ok que de un error y ponga a nil la vm y salga
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
	// Falta devolver cual es el ID de la nueva ademas de eso tendras que hacer un PUT para modificar los
	// parametros de la instancia nueva. entre ellos el procesador la memoria y la network
	//
	// Después de esto tendrás que crear una network y asignarla a la instancia nueva.
	return &vm, err
}

// Method ReadVM return the object MyVm with the ID indicate in i.
// Input: i: string with the ID of the VM, Return: pointer at the MyVm object
// and error variable with the error if occurr
func (c *Client) ReadVM(i string) (*MyVm, error) {
	var vms []MyVm
	var vm MyVm
	// La única manera de poder conseguir el path de la instancia
	// es pidiendo a la API de VMW todas las instancias
	responseBody, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(responseBody).Decode(&vms)
	if err != nil {
		return nil, err
	}

	for id, value := range vms {
		if value.IdVM == i {
			vm = vms[id]
		}
	}

	file, err := os.Open(vm.Path)
	if err != nil {
		log.Fatalf("Failed opening file %s, please make sure the config file exists", err)
	}

	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:File object %#v\n", file)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var vardenomination, vardescription string
	for scanner.Scan() {
		temp := strings.SplitN(scanner.Text(), " = ", 2)
		vtemp := strings.ToLower(temp[0])
		if vtemp == "displayname" {
			vardenomination = strings.TrimSpace(temp[1])
			vardenomination = strings.TrimSuffix(vardenomination, "\"")
			vardenomination = strings.TrimPrefix(vardenomination, "\"")
		}
		if vtemp == "annotation" {
			vardescription = strings.TrimSpace(temp[1])
			vardescription = strings.TrimSuffix(vardescription, "\"")
			vardescription = strings.TrimPrefix(vardescription, "\"")
		}
	}
	vm.Denomination = vardenomination
	vm.Description = vardescription
	file.Close()

	responseBody, err = c.httpRequest("vms/"+i, "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:Body of VM %#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}
	responseBody, err = c.httpRequest("vms/"+i+"/power", "GET", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:Body of power %#v\n", responseBody)
	err = json.NewDecoder(responseBody).Decode(&vm)
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: ReadVM Obj:VM %#v\n", vm)
	return &vm, nil
}

// Function to update a VM in VmWare Worstation Input:
// i: string with the ID of the VM to update, d: string with the denomination of the VM
func (c *Client) UpdateVM(i string) (*MyVm, error) {
	var vm MyVm
	responseBody, err := c.httpRequest("vms/"+i, "PUT", bytes.Buffer{})
	if err != nil {
		return nil, err
	}
	log.Printf("[WSAPICLI] Fi: wsapivm.go Fu: UpdateVM Obj:%#v\n", responseBody)

	return &vm, err
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
