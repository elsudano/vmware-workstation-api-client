package wsapiutils

import (
	"os"

	"github.com/elsudano/govmx"
	"github.com/rs/zerolog/log"
)

func New(vmx *govmx.VirtualMachine) VMFile {
	return &VMStructure{myvm: vmx}
}

// GetVMFromFile - With this function we can obtain a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Inputs:
// f: string, the complete path of the vxm file that we want to read
// Outputs:
// string, vmx.VirtualMachine structure, and error if you obtain some error in the function
func (vmf *VMStructure) GetVMFromFile(f string) error {
	log.Info().Msgf("We are trying to read the file: %#v", f)
	data, err := os.ReadFile(f)
	if err != nil {
		log.Error().Err(err).Msg("Please make sure the config file exists.")
		return err
	}
	log.Debug().Msgf("After read the file: %#v", data)
	err = govmx.Unmarshal(data, vmf.myvm)
	if err != nil {
		log.Error().Err(err).Msg("Error trying to Unmarshal the data.")
		return err
	}
	log.Debug().Msgf("After Unmarshal the data: %#v", vmf.myvm)
	log.Info().Msg("We have read the VM file and load the data in the vmx object.")
	return nil
}

// SetVMToFile - With this function we can save a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Inputs:
// f: (string), File where we want to save the VM
// p: (string), Chain with the parameter we want to change
// v: (string), Chain with the value of the parameter that we want change
// Output:
// error if you obtain some error in the function
func (vmf *VMStructure) SetVMToFile(f string) error {
	log.Info().Msgf("We will write in this file: %#v.", f)
	data, err := govmx.Marshal(vmf.myvm)
	if err != nil {
		log.Error().Err(err).Msg("Error trying to Unmarshal the data.")
		return err
	}
	log.Debug().Msgf("After Marshaling the data VM is: %#v", vmf.myvm)
	err = os.WriteFile(f, data, 0644)
	if err != nil {
		log.Error().Err(err).Msgf("Error writing in file %#v, please make sure the file exists.", f)
		return err
	}
	log.Debug().Msgf("We have saved the VM %#v in the file: %#v", vmf.myvm, f)
	log.Info().Msg("We have saved the VM in a file.")
	return err
}

// GetAnnotation - With this function we can obtain the value of the description of VM
// Input:
// f: string, the complete path of the vxm file that we want to read
// Output:
// string, Value of the Annotation field of the VM.
// error if you obtain some error in the function.
func (vmf *VMStructure) GetAnnotation(f string) (string, error) {
	err := vmf.GetVMFromFile(f)
	if err != nil {
		log.Error().Err(err).Msg("Failure to obtain the value of the Description.")
		return "", err
	}
	return vmf.myvm.Annotation, nil
}

// SetAnnotation - With this function we can set the value of the description of VM
// Input:
// f: (string) The complete path of the vxm file that we want to write
// v: (string) The value of Annotation field
// Output:
// error if you obtain some error in the function
func (vmf *VMStructure) SetAnnotation(f string, v string) error {
	vmf.myvm.Annotation = v
	err := vmf.SetVMToFile(f)
	if err != nil {
		log.Error().Err(err).Msgf("We haven't be able to save the vmx data in the file %#v", f)
		return err
	}
	return nil
}

// GetDisplayName - With this function we can obtain the name of VM
// Input:
// f: string, the complete path of the vxm file that we want to read
// Output:
// string, Value of the Denomination field of the VM.
// error if you obtain some error in the function.
func (vmf *VMStructure) GetDisplayName(f string) (string, error) {
	err := vmf.GetVMFromFile(f)
	if err != nil {
		log.Error().Err(err).Msg("Failure to obtain the value of the Denomination.")
		return "", err
	}
	return vmf.myvm.DisplayName, nil
}

// SetDisplayName - With this function we can set the value of the denomination of VM
// Input:
// f: string, the complete path of the vxm file that we want to write.
// v: string with the value of Denomination field, WARNING this function don't change the PATH
// Output:
// error if you obtain some error in the function.
func (vmf *VMStructure) SetDisplayName(f string, v string) error {
	vmf.myvm.DisplayName = v
	// Here you will need to change the folder and the file name
	// when we changed the displayName
	err := vmf.SetVMToFile(f)
	if err != nil {
		log.Error().Err(err).Msgf("We haven't be able to save the vmx data in the file %#v", f)
		return err
	}
	return nil
}

// SetDenominationDescription With this function you can setting the Denomination and Description of the VM.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm.
// Input:
// f: (string) The complete path of the vxm file that we want to modify.
// n: (string) The value of denomination.
// d: (string) The value of description.
// Output:
// err: variable with error if occur
func (vmf *VMStructure) SetDenominationDescription(f string, n string, d string) error {
	log.Info().Msgf("The new values for Denomination %#v, and Description. %#v", n, d)
	err := vmf.SetDisplayName(f, n)
	if err != nil {
		log.Error().Err(err).Msg("We can't change the Denomination.")
		return err
	}
	err = vmf.SetAnnotation(f, d)
	if err != nil {
		log.Error().Err(err).Msg("We can't change the Description.")
		return err
	}
	return nil
}
