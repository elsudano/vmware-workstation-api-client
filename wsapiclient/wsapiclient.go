package wsapiclient

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/elsudano/vmware-workstation-api-client/httpclient"
	"github.com/elsudano/vmware-workstation-api-client/wsapinet"
	"github.com/elsudano/vmware-workstation-api-client/wsapivm"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	LibraryVersion string = "2.7.45"
)

type WSAPIService interface {
	ConfigLog(lvl string, mode string)
	ConfigApiClient(a string, u string, p string, i bool, d string) error
	GetAllVMs() ([]wsapivm.MyVm, error)
	LoadVM(i string) (*wsapivm.MyVm, error)
	LoadVMbyName(n string) (*wsapivm.MyVm, error)
	CreateVM(pid string, n string, d string, p int, m int) (*wsapivm.MyVm, error)
	UpdateVM(vm *wsapivm.MyVm, n string, d string, p int, m int, s string) error
	RegisterVM(vm *wsapivm.MyVm) error
	DeleteVM(vm *wsapivm.MyVm) error
}

type WSAPIClient struct {
	Caller     *httpclient.HTTPClient
	VMService  wsapivm.VMService
	NETService wsapinet.NETService
}

func New() WSAPIService {
	myclient, err := httpclient.New()
	if err != nil {
		log.Error().Err(err).Msg("We can't create the Client.")
		return nil
	}
	return &WSAPIClient{
		Caller:     myclient,
		VMService:  wsapivm.New(myclient),
		NETService: wsapinet.New(myclient),
	}
}

// ConfigLog method change the behavior that how to handle the logging on our API
// Inputs:
// lvl: (strings) Which will be the level bu default that we want in console
// mode: (string) The format we want:
// (defaut) JSON: in stderr,
// FILE json format in debug.log file,
// CONSOLE stdout with color in json format,
// ERROR stderr in json format
// HR (Human Readable) in stdout
func (wsapi *WSAPIClient) ConfigLog(lvl string, mode string) {
	lvl = strings.ToUpper(lvl)
	switch lvl {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
	// Global Settings https://github.com/rs/zerolog?tab=readme-ov-file#global-settings
	zerolog.TimeFieldFormat = time.RFC3339 // zerolog.TimeFormatUnix zerolog.TimeFormatUnixMs, zerolog.TimeFormatUnixMicro

	// Customized Fields Name https://github.com/rs/zerolog?tab=readme-ov-file#customize-automatic-field-names
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	// To trace the errors https://github.com/rs/zerolog?tab=readme-ov-file#add-file-and-line-number-to-log
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()

	// Formatting https://github.com/rs/zerolog?tab=readme-ov-file#pretty-logging
	switch strings.ToUpper(mode) {
	case "FILE":
		file, err := os.OpenFile(
			"debug.log",
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		ConsoleWriter := zerolog.ConsoleWriter{Out: file, NoColor: false}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "CONSOLE":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		ConsoleWriter.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}
		ConsoleWriter.PartsExclude = []string{
			zerolog.TimestampFieldName,
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "ERROR":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-6s][WSAPICLI]", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	case "HR":
		ConsoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}
		ConsoleWriter.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("[%-6s][WSAPICLI]", i))
		}
		ConsoleWriter.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		ConsoleWriter.PartsExclude = []string{
			zerolog.TimestampFieldName,
		}
		log.Logger = log.With().Logger().Output(ConsoleWriter)
	default:
		log.Logger = log.With().Logger().Output(os.Stderr)
	}
}

// ConfigApiClient method return a pointer of Client of API but now it's configure
// Inputs:
// waspi: (WSAPIClient) client with all the necessary data to make a call.
// a: (string) address of URL to server of API.
// u: (string) user for to authenticate.
// p: (string) password of user.
// i: (bool) Insecure flag to http or https.
// d: (string) debug mode
func (wsapi *WSAPIClient) ConfigApiClient(a string, u string, p string, i bool, d string) error {
	return wsapi.Caller.ConfigClient(a, u, p, i, d)
}

// GetAllVMs Method return array of MyVm and a error variable if occurr some problem
// Outputs:
// []MyVm list of all VMs that we have in VmWare Workstation
// (error) variable with the error if occurr
func (wsapi *WSAPIClient) GetAllVMs() ([]wsapivm.MyVm, error) {
	return wsapi.VMService.GetAllVMs()
}

// CreateVM method to create a new VM in VmWare Worstation
// Input:
// pid: (string) with the ID of the Parent VM,
// n: string with the denomination of the VM,
// d: string with the description of VM
// p: int with the number of processors in the VM
// m: int with the number of memory in the VM
func (wsapi *WSAPIClient) CreateVM(pid string, n string, d string, p int, m int) (*wsapivm.MyVm, error) {
	return nil, nil
}

// LoadVM method return the object MyVm with the ID indicate in i.
// Inputs:
// i: (string) String with the ID of the VM
// Outputs:
// (pointer) Pointer at the MyVm object
// (error) variable with the error if occurr
func (wsapi *WSAPIClient) LoadVM(i string) (*wsapivm.MyVm, error) { return nil, nil }

// LoadVMbyName method return the object MyVm with the Name indicate in n.
// Inputs:
// n: (string) String with the Name of the VM
// Outputs:
// (pointer) Pointer at the MyVm object
// (error) variable with the error if occurr
func (wsapi *WSAPIClient) LoadVMbyName(n string) (*wsapivm.MyVm, error) { return nil, nil }

// UpdateVM method to update a VM in VmWare Worstation
// Input:
// vm (*MyVm) The VM that we want to update
// n: string with the denomination of VM
// d: string with the description of the VM, p: int with the number of processors
// m: int with the size of memory
// s: Power State desired, choose between on, off, reset, (nil no change)
// Output:
// pointer at the MyVm object
// and error variable with the error if occurr
func (wsapi *WSAPIClient) UpdateVM(vm *wsapivm.MyVm, n string, d string, p int, m int, s string) error {
	return nil
}

// RegisterVM method to register a new VM in VmWare Worstation GUI:
// Input:
// c: (*wsapiclient.Client) The client to make the call.
// vm: (*wsapivm.MyVM) The VM object that we want to delete.
// Output:
// error: (error) The possible error that you will have.
func (wsapi *WSAPIClient) RegisterVM(vm *wsapivm.MyVm) error {
	return nil
}

// DeleteVM method to delete a VM in VmWare Worstation
// Input:
// c: (*wsapiclient.Client) The client to make the call.
// vm: (*wsapivm.MyVM) The VM object that we want to delete.
// Output:
// error: (error) The possible error that you will have.
func (wsapi *WSAPIClient) DeleteVM(vm *wsapivm.MyVm) error {
	return nil
}
