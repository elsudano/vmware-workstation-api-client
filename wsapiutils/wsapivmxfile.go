package wsapiutils

import (
	"log"
	"os"

	vmx "github.com/johlandabee/govmx"
)

// GetVMFromFile - With this function we can obtain a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Inputs:
// p: string, the complete path of the vxm file that we want to read
// Outputs:
// string, vmx.VirtualMachine structure, and error if you obtain some error in the function
func GetVMFromFile(p string) (vmx.VirtualMachine, error) {
	vm := new(vmx.VirtualMachine)
	data, err := os.ReadFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Message: Failed %s, please make sure the config file exists", err)
		return *vm, err
	}

	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Obj: Data File %#v\n", string(data))
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Obj: %#v", err)
		return *vm, err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: GetVMFromFile Obj: VM  %#v\n", vm)
	return *vm, nil
}

// SetVMToFile - With this function we can save a vmx.VirtualMachine structure
// with all the possible values that we have in the file.
// Inputs:
// p: string, with the parameter we want to change
// Output:
// error if you obtain some error in the function
func SetVMToFile(vm vmx.VirtualMachine, p string) error {
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Message: parameters %#v, %#v", vm, p)
	data, err := vmx.Marshal(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Message: Failed to save the VMX structure in memory %s", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Obj: Data after read vm %#v\n", string(data))
	err = os.WriteFile(p, data, 0644)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetVMToFile Message: Failed writing in file %s, please make sure the config file exists", err)
		return err
	}
	return err
}

// GetAnnotation - With this function we can obtain the value of the description of VM
// Input:
// p: string, the complete path of the vxm file that we want to read
// Output:
// string, Value of the Annotation field of the VM, error if you obtain some error in the fuction
func GetAnnotation(p string) (string, error) {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetAnnotation Message: Failure to obtain the value of the Description %s", err)
		return "", err
	}
	return vm.Annotation, nil
}

// SetAnnotation - With this function we can set the value of the description of VM
// Input: p: string, the complete path of the vxm file that we want to read
// v: string with the value of Annotation field
// Output: error if you obtain some error in the fuction
func SetAnnotation(p string, v string) error {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetAnnotation Message: We can't obtain the vmx object %s", err)
		return err
	}
	vm.Annotation = v
	err = SetVMToFile(vm, p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetAnnotation Message: We haven't be able to save the structure in the file %s", err)
		return err
	}
	return nil
}

// GetDisplayName - With this function we can obtain the value of the name of VM
// Input: p: string, the complete path of the vxm file that we want to read
// Output: string, Value of the Denomination field of the VM, error if you obtain some error in the fuction
func GetDisplayName(p string) (string, error) {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: GetDisplayName Message: Failure to obtain the value of the Denomination %s", err)
		return "", err
	}
	return vm.DisplayName, nil
}

// SetDisplayName - With this function we can set the value of the denomination of VM
// Input: p: string, the complete path of the vxm file that we want to read
// v: string with the value of Denomination field, WARNING this function don't change teh PATH
// Output: error if you obtain some error in the fuction
func SetDisplayName(p string, v string) error {
	vm, err := GetVMFromFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetAnnotation Message: We can't obtain the vmx object %s", err)
		return err
	}
	vm.DisplayName = v
	err = SetVMToFile(vm, p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetAnnotation Message: We haven't be able to save the structure in the file %s", err)
		return err
	}
	return nil
}

// SetNameDescription With this function you can setting the Denomination and Description of the VM.
// this information is in the vmx file of the machine for that you need know
// which is the file of the vm. Input: p: string with the complete path of the file,
// n: string with the denomination, d: string with the description err: variable with error if occur
func SetNameDescription(p string, n string, d string) error {
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: parameters %#v, %#v, %#v", p, n, d)
	data, err := os.ReadFile(p)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed opening file %s, please make sure the config file exists", err)
		return err
	}

	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: File object %#v\n", string(data))

	vm := new(vmx.VirtualMachine)
	err = vmx.Unmarshal(data, vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: %#v", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: VM %#v\n", vm)

	vm.DisplayName = n
	vm.Annotation = d
	data, err = vmx.Marshal(vm)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed to save the VMX structure in memory %s", err)
		return err
	}
	log.Printf("[DEBUG][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Obj: Data File %#v\n", string(data))
	err = os.WriteFile(p, data, 0644)
	if err != nil {
		log.Printf("[ERROR][WSAPICLI] Fi: wsapitools.go Fu: SetNameDescription Message: Failed writing in file %s, please make sure the config file exists", err)
		return err
	}
	// en este punto tambien tienes que cambiar el nombre del fihero cuando se cambia la denominacion
	return err
}
