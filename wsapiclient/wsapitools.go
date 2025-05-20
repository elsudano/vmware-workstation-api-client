package wsapiclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
)

// GetVM Auxiliar function to get the data of the VM and don't repeat code
// Input:
// c: (pointer) pointer at the client of the API server
// i: (string) string with the ID yo VM
// Outputs:
// vm: (pointer) pointer to the VM that we are handeling
// err: (error) If we will have some error we can handle it here.
func (c *Client) GetVM(i string) (*MyVm, error) {
	var vms []MyVm
	var vm MyVm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't allow see this the another way
	// --------- This Block read the path and the ID of the vm in order to load in the function --------- {{{
	response, vmerror, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: The request at the server API failed %s", err)
		return nil, err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&vms)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetVM Message: I can't read the json structure %s", err)
			return nil, err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo Obj: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
		}
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetVM Obj: VM %#v\n", vm)
	return &vm, nil
}

// GetBasicInfo Auxiliar function in charge of getting de Basic Information
// Inputs:
// c: (pointer) Pointer at the client of the API server
// vm: (MyVm) The VM that we want to know the Memory and CPU info
// Outputs:
// err: (error) If we will have some error we can handle it here.
func (c *Client) GetBasicInfo(vm *MyVm) error {
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM, "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo M: Request Error Getting Information %#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo M: Response Raw Information %#v\n", response)
		err = json.NewDecoder(response).Decode(&vm)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo M: Response Error Getting Information %#v\n", err)
			return err
		}
	case 110:
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo Code %d M: %s", vmerror.Code, vmerror.Message)
		return errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetBasicInfo Obj: VM %#v\n", vm)
	return nil
}

// GetDenominationDescription Axiliar function in charge about the getting the
// description and Denomination of the VM and set in our structure.
// Inputs:
// c: (pointer) Pointer at the client of the API server
// vm: (MyVm) that's the VM where we want use the information.
// Outputs:
// err: (error) If we will have some error we can handle it here.
func (c *Client) GetDenominationDescription(vm *MyVm) error {
	var param ParamPayload
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/params/displayName", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Request Error Getting Denomination %#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Response Getting Denomination %#v\n", response)
		err = json.NewDecoder(response).Decode(&param)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Error Decoding Denomination %#v\n", err)
			return err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetDenominationDescription M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	vm.Denomination = param.Value
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/params/annotation", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Request Error Getting Description %#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Response Getting Description %#v\n", response)
		err = json.NewDecoder(response).Decode(&param)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription M: Error Decoding Description %#v\n", err)
			return err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: GetDenominationDescription M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	vm.Description = param.Value
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetDenominationDescription Obj: VM %#v\n", vm)
	return nil
}

// GetPowerStatus Auxiliar function in charge to get the current Power Status
// Inputs:
// c: (pointer) at the client of the API server
// vm: (MyVm) The VM that we want to know the Power Status
// Outputs:
// err: (error) If we will have some error we can handle it here.
func (c *Client) GetPowerStatus(vm *MyVm) error {
	var power_state_payload PowerStatePayload
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/power", "GET", bytes.Buffer{})
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetPowerStatus Obj: Request Error Power Status %#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetPowerStatus Obj: Response Power Status %#v\n", response)
		err = json.NewDecoder(response).Decode(&power_state_payload)
		vm.PowerStatus = PowerStateConversor(power_state_payload.Value)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetPowerStatus Obj: Decoding Power Status %#v\n", err)
			return err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetPowerStatus Obj: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetPowerStatus Obj: VM %#v\n", vm)
	return nil
}

// PowerSwitch method that permit you change the state of the instance, so you will change
// from power-off to power-on the state of the instance.
// Inputs:
// c: (pointer) Pointer at the client of the API server
// vm: (MyVm) the VM object that we want to change the Power Status,
// s: (string) String with the state that will want between on, off, reset
// Outputs:
// err: (error) If we will have some error we can handle it here.
func (c *Client) PowerSwitch(vm *MyVm, s string) error {
	var power_state_payload PowerStatePayload
	requestBody := bytes.NewBufferString(s)
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch Obj: Request option %#v\n", requestBody.String())
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/power", "PUT", *requestBody)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch M: Response RAW %#v\n", err)
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&power_state_payload)
		vm.PowerStatus = PowerStateConversor(power_state_payload.Value)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapivm.go Fu: PowerSwitch M: Response Body RAW %#v, %#v\n", err, response)
			return err
		}
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: PowerSwitch M: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: PowerSwitch Obj: VM %#v\n", vm)
	return nil
}

// PowerStateConversor We have to create this method, because the API of th VMWare Workstation
// change the values of the Power State of the instance, I mean, If I send "on" the API change
// the value for powerOn, and obviusly that is a big problem
// Inputs:
// ops: (string) The original Power State, the string that the API of VmWare Workstation give us
// Outputs:
// s: (string) The normalized string
func PowerStateConversor(ops string) (s string) {
	switch ops {
	case "poweredOn":
		return "on"
	case "poweringOn":
		return "on"
	case "poweredOff":
		return "off"
	case "poweringOff":
		return "off"
	default:
		return "Invalid Power State"
	}
}

// SetParameter With this function you can set the value of the parameter.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm.
// Input:
// c: (pointer) Pointer at the client of the API server
// vm: (MyVm) is the VM object that we will changes,
// p: (string) String with the name or param to set,
// v: (string) String with the value of param err: variable with error if occur
// Outputs:
// err: (error) If we will have some error we can handle it here.
func (c *Client) SetParameter(vm *MyVm, p string, v string) error {
	var param ParamPayload
	param.Name = p
	param.Value = v
	request, err := json.Marshal(param)
	// request, err := json.Marshal(map[string]string{
	// 	"name":  p,
	// 	"value": v,
	// })
	if err != nil {
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:%#v\n", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request %#v\n", request)
	requestBody := new(bytes.Buffer)
	requestBody.Write(request)
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Request Body %#v\n", requestBody.String())
	response, vmerror, err := c.httpRequest("/vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	if err != nil {
		return err
	}
	switch vmerror.Code {
	case 0:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:response raw %#v\n", response)
		responseBody := new(bytes.Buffer)
		_, err = responseBody.ReadFrom(response)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Response Error %#v\n", err)
			return err
		}
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj:Response Body %#v\n", responseBody.String())
	default:
		log.Printf("[DEBUG][WSAPICLI] Fi: wsapinet.go Fu: SetParameter Obj: Output Code %d and Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetParameter Obj: VM %#v\n", vm)
	return nil
}

// InitialData is a extra function to fill the fields that we need to use
// to make the tests in our API.
// Inputs:
// f: (string) is the file where we have the configuration.
// Outputs:
// url: (string) That will be the URL of our API endpoint
// user: (string) That will be the User of our API
// pass: (string) That will be the Password of our API
// parentid: (string) That will be the Parent ID of our VM
// insecure: (bool) If our API works with HTTP we will set true here
// debug: (bool) If we need troubleshot our API we will set true here
// error: (error) When the function catch some error print it here
func InitialData(f string) (string, string, string, string, bool, bool, error) {
	var user, pass, url, parentid string
	var insecure, debug = false, false
	fileInfo, err := os.Stat(f)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: InitialData M: Error while we trying to check the file: %#v", err)
		return "", "", "", "", false, false, err
	}
	if fileInfo.Mode().IsDir() {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: InitialData M: It is a directory, please select a config file")
		return "", "", "", "", false, false, nil
	} else if fileInfo.Mode().IsRegular() {
		file, err := os.Open(f)
		if err != nil {
			log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: InitialData M: Failed opening file %s, please make sure the config file exists", err)
			return "", "", "", "", false, false, err
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			temp := strings.SplitN(scanner.Text(), ":", 2)
			key := strings.ToLower(temp[0])
			if key == "user" {
				user = strings.TrimSpace(temp[1])
			}
			if key == "password" {
				pass = strings.TrimSpace(temp[1])
			}
			if key == "baseurl" {
				url = strings.TrimSpace(temp[1])
			}
			if key == "parentid" {
				parentid = strings.TrimSpace(temp[1])
			}
			if key == "insecure" && strings.TrimSpace(temp[1]) == "true" {
				insecure = true
			}
			if key == "debug" && strings.TrimSpace(temp[1]) == "true" {
				debug = true
			}
		}
	} else {
		log.Println("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: InitialData M: Something was wrong, please try again")
		return "", "", "", "", false, false, nil
	}
	return url, user, pass, parentid, insecure, debug, nil
}
