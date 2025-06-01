package wsapiclient

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// GetVM Auxiliar function to get the data of the VM and don't repeat code
// Input:
// c: (pointer) pointer at the client of the API server
// i: (string) string with the ID yo VM
// Outputs:
// vm: (pointer) pointer to the VM that we are handeling
// err: (error) If we will have some error we can handle it here.
func (c *Client) GetVM(i string) (*MyVm, error) {
	log.Info().Msgf("The VM Id value is: %#v", i)
	var vms []MyVm
	var vm MyVm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't allow see this the another way
	// --------- This Block read the path and the ID of the vm in order to load in the function --------- {{{
	response, vmerror, err := c.httpRequest("vms", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We can't made the API call.")
		return nil, err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&vms)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return nil, err
		}
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return nil, errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("List of VMs: %#v", vms)
	for tempvm, value := range vms {
		if value.IdVM == i {
			vm = vms[tempvm]
			// poner un break cuando la encuentre
		}
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the ID and Path values.")
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
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&vm)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
	case 110:
		err = errors.New("Code:" + strconv.Itoa(vmerror.Code) + " Msg:" + vmerror.Message)
		log.Error().Msgf("Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return err
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the Processor and Memory values.")
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
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&param)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	vm.Denomination = param.Value
	response, vmerror, err = c.httpRequest("vms/"+vm.IdVM+"/params/annotation", "GET", bytes.Buffer{})
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&param)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	vm.Description = param.Value
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the Denomination and Description values.")
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
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&power_state_payload)
		vm.PowerStatus = PowerStateConversor(power_state_payload.Value)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the Power State value.")
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
	log.Debug().Msgf("The state that we want is: %#v", s)
	response, vmerror, err := c.httpRequest("vms/"+vm.IdVM+"/power", "PUT", *requestBody)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't complete the API call.")
		return err
	}
	switch vmerror.Code {
	case 0:
		err = json.NewDecoder(response).Decode(&power_state_payload)
		vm.PowerStatus = PowerStateConversor(power_state_payload.Value)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return err
		}
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have changed the Power State.")
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
	log.Debug().Msgf("Power State RAW: %#v", ops)
	log.Info().Msg("We have converted the Power State.")
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
	if err != nil {
		log.Error().Err(err).Msgf("Trying to encode this request: %#v", param)
		return err
	}
	requestBody := new(bytes.Buffer)
	requestBody.Write(request)
	log.Debug().Msgf("Request Human Readable: %#v", requestBody.String())
	response, vmerror, err := c.httpRequest("/vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
	if err != nil {
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
	default:
		log.Error().Msgf("We haven't handled this error Code: %d Message: %s", vmerror.Code, vmerror.Message)
		return errors.New(strconv.Itoa(vmerror.Code) + "," + vmerror.Message)
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msgf("We have defined new value in parameter: %#v", p)
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
func InitialData(f string) (string, string, string, string, bool, string, error) {
	var user, pass, url, debug, parentid string
	var insecure = false
	fileInfo, err := os.Stat(f)
	if err != nil {
		log.Error().Err(err).Msg("While we trying to check the file.")
		return "", "", "", "", false, "", err
	}
	if fileInfo.Mode().IsDir() {
		log.Error().Err(err).Msg("It is a directory, please select a config file.")
		return "", "", "", "", false, "", nil
	} else if fileInfo.Mode().IsRegular() {
		file, err := os.Open(f)
		if err != nil {
			log.Error().Err(err).Msg("Failed opening file, please make sure the config file exists.")
			return "", "", "", "", false, "", err
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
			if key == "debug" {
				debug = strings.TrimSpace(temp[1])
			}
		}
	} else {
		log.Error().Msgf("We haven't handled this error, something was wrong, please try again")
		return "", "", "", "", false, "", nil
	}
	return url, user, pass, parentid, insecure, debug, nil
}
