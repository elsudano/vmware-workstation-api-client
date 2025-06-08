package wsapiutils

import vmx "github.com/johlandabee/govmx"

// Interface with all the methods that we can use to talk with the API of VmWare Workstation Pro
type VMFile interface {
	GetVMFromFile(p string) error
	SetVMToFile(f string) error
	GetAnnotation(f string) (string, error)
	SetAnnotation(f string, v string) error
	GetDisplayName(f string) (string, error)
	SetDisplayName(f string, v string) error
	SetDenominationDescription(f string, n string, d string) error
}

// That's the abstract object that how we see our VM's
type VMStructure struct {
	myvm *vmx.VirtualMachine
}
