package main

import (
	"fmt"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
)

func main() {
	client, err := wsapiclient.New()
	// client.SwitchDebug()
	AllVMs, err := client.GetAllVMs()
	if err != nil {
		panic(err)
	}
	for _, value := range AllVMs {
		fmt.Printf("ID: %v\nPath: %v\nDescription: %s\nPower Status: %s\nProcessor: %d\nMemory: %d \n\n",
			value.IdVM,
			value.Path,
			value.Description,
			value.PowerStatus,
			value.CPU.Processors,
			value.Memory)
	}
	VM := client.GetVM(AllVMs[0].IdVM)
	fmt.Printf("VM Id: %v\nVM Processors: %d\nVM Memory: %d\nVM Status: %v\n",
		VM.IdVM, VM.CPU.Processors, VM.Memory, VM.PowerStatus)
}
