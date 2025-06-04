package wsapivm

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/rs/zerolog/log"
)

// GetVM Auxiliar function to get the data of the VM and don't repeat code
// Input:
// c: (*wsapiclient.Client) pointer at the client of the API server.
// i: (string) string with the ID yo VM
// Outputs:
// vm: (*wsapivm.MyVm) pointer to the VM that we are handeling.
// err: (error) If we will have some error we can handle it here.
func GetVM(vmc *httpclient.HTTPClient, i string) (*MyVm, error) {
	log.Info().Msgf("The VM Id value is: %#v", i)
	var vms []MyVm
	var vm MyVm
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't allow see this the another way
	// --------- This Block read the path and the ID of the vm in order to load in the function --------- {{{
	response, vmerror, err := vmc.ApiCall("vms", "GET", bytes.Buffer{})
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
			break
		}
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the ID and Path values.")
	return &vm, nil
}

// GetVMbyName Auxiliary function to get the data of the VM and don't repeat code
// Input:
// vmc: (*httpclient.HTTPClient) pointer at the client of the API server.
// n: (string) The name of the VM that we want to get.
// Outputs:
// vm: (*wsapivm.MyVm) pointer to the VM that we are handeling.
// err: (error) If we will have some error we can handle it here.
func GetVMbyName(vmc *httpclient.HTTPClient, n string) (*MyVm, error) {
	log.Info().Msgf("The VM name value is: %#v", n)
	var vms []MyVm
	var vm MyVm
	var param ParamPayload
	// If you want see the path of the VM it's necessary getting all VMs
	// because the API of VmWare Workstation doesn't allow see this the another way
	// --------- This Block read the path and the ID of the vm in order to load in the function --------- {{{
	response, vmerror, err := vmc.ApiCall("vms", "GET", bytes.Buffer{})
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
		response, _, err = vmc.ApiCall("vms/"+value.IdVM+"/params/displayName", "GET", bytes.Buffer{})
		if err != nil {
			log.Error().Err(err).Msg("We couldn't complete the API call.")
			return nil, err
		}
		err = json.NewDecoder(response).Decode(&param)
		if err != nil {
			log.Error().Err(err).Msg("The response JSON is malformed.")
			return nil, err
		}
		if param.Value == n {
			vm = vms[tempvm]
			break
		}
	}
	log.Debug().Msgf("VM: %#v", vm)
	log.Info().Msg("We have loaded the ID and Path values.")
	return &vm, nil
}

// GetAllExtraParameters Auxiliary function to get all the Extra parameters
// we have created this function in order not repeat the same code in both
// functions LoadVM and LoadVMbyName.
// Inputs:
// vm: (*wsapivm.MyVm) That's will be the pointer at our vm that we want fill.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// Outputs:
// err: (error) If we will have some error we can handle it here.
func GetAllExtraParameters(vmc *httpclient.HTTPClient, vm *MyVm) error {
	err := GetBasicInfo(vmc, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Process and Memory.")
		return err
	}
	err = GetDenominationDescription(vmc, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Denomination and Description.")
		return err
	}
	err = GetPowerStatus(vmc, vm)
	if err != nil {
		log.Error().Err(err).Msg("We couldn't show the Power Status.")
		return err
	}
	// if vm.PowerStatus == "on" {
	// 	err = wsapinet.GetInfoNics(vmc, vm.IdVM)
	// 	if err != nil {
	// 		log.Error().Err(err).Msg("We couldn't show the Network information.")
	// 		return err
	// 	}
	// }
	return nil
}

// GetBasicInfo Auxiliary function in charge of getting de Basic Information
// Inputs:
// vm: (*wsapivm.MyVm) The VM that we want to know the Memory and CPU info.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// Outputs:
// err: (error) If we will have some error we can handle it here.
func GetBasicInfo(vmc *httpclient.HTTPClient, vm *MyVm) error {
	response, vmerror, err := vmc.ApiCall("vms/"+vm.IdVM, "GET", bytes.Buffer{})
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

// GetDenominationDescription Auxiliary function in charge about the getting the
// description and Denomination of the VM and set in our structure.
// Inputs:
// vm: (*wsapivm.MyVm) The VM that we want to know the Denomination and Description info.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// Outputs:
// err: (error) If we will have some error we can handle it here.
func GetDenominationDescription(vmc *httpclient.HTTPClient, vm *MyVm) error {
	var param ParamPayload
	response, vmerror, err := vmc.ApiCall("vms/"+vm.IdVM+"/params/displayName", "GET", bytes.Buffer{})
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
	response, vmerror, err = vmc.ApiCall("vms/"+vm.IdVM+"/params/annotation", "GET", bytes.Buffer{})
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

// GetPowerStatus Auxiliary function in charge to get the current Power Status
// Inputs:
// vm: (*wsapivm.MyVm) The VM that we want to know the Power Status info.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// Outputs:
// err: (error) If we will have some error we can handle it here.
func GetPowerStatus(vmc *httpclient.HTTPClient, vm *MyVm) error {
	var power_state_payload PowerStatePayload
	response, vmerror, err := vmc.ApiCall("vms/"+vm.IdVM+"/power", "GET", bytes.Buffer{})
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
// vm: (*wsapivm.MyVm) The VM that we want to know the Denomination and Description info.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// s: (string) String with the state that will want between on, off, reset
// Outputs:
// err: (error) If we will have some error we can handle it here.
func PowerSwitch(vmc *httpclient.HTTPClient, vm *MyVm, s string) error {
	var power_state_payload PowerStatePayload
	requestBody := bytes.NewBufferString(s)
	log.Debug().Msgf("The state that we want is: %#v", s)
	response, vmerror, err := vmc.ApiCall("vms/"+vm.IdVM+"/power", "PUT", *requestBody)
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
// the value for powerOn, and obviously that is a big problem
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
// Inputs:
// vm: (*wsapivm.MyVm) The VM that we want to know the Denomination and Description info.
// c: (*wsapiclient.Client) pointer at the client of the API server.
// p: (string) String with the name or param to set,
// v: (string) String with the value of param err: variable with error if occur
// Outputs:
// err: (error) If we will have some error we can handle it here.
func SetParameter(vmc *httpclient.HTTPClient, vm *MyVm, p string, v string) error {
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
	response, vmerror, err := vmc.ApiCall("/vms/"+vm.IdVM+"/configparams", "PUT", *requestBody)
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
