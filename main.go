package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
)

func main() {
	file, err := os.Open("config.ini")
	if err != nil {
		log.Fatalf("Failed opening file %s, please make sure the config file exists", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var varuser, varpass, varurl, varparentid string
	var varinsecure, vardebug bool
	for scanner.Scan() {
		temp := strings.SplitN(scanner.Text(), ":", 2)
		vtemp := strings.ToLower(temp[0])
		vtemp2 := temp[1]
		fmt.Printf("Temp: %s, Key: %s, Value: %s\n", temp, vtemp, vtemp2)
		if vtemp == "user" {
			varuser = strings.TrimSpace(temp[1])
		}
		if vtemp == "password" {
			varpass = strings.TrimSpace(temp[1])
		}
		if vtemp == "baseurl" {
			varurl = strings.TrimSpace(temp[1])
		}
		if vtemp == "parentid" {
			varparentid = strings.TrimSpace(temp[1])
		}
		if vtemp == "insecure" {
			if strings.TrimSpace(temp[1]) == "true" {
				varinsecure = true
			} else {
				varinsecure = false
			}
		}
		if vtemp == "debug" {
			if strings.TrimSpace(temp[1]) == "true" {
				vardebug = true
			} else {
				vardebug = false
			}
		}
	}
	fmt.Printf("varuser: %s, varpass: %s, varurl: %s, varparentid: %s, varinsecure: %v, vardebug: %v\n", varuser, varpass, varurl, varparentid, varinsecure, vardebug)
	file.Close()

	// First option we will can create a client with the default values {{{
	// client, _ := wsapiclient.New()
	// client.ConfigCli(varurl, varuser, varpass, varinsecure, vardebug)
	// fmt.Printf("Parent ID: %s\n", varparentid)
	// }}}

	// Second option we will setting the values in the creation moment {{{
	client, err := wsapiclient.NewClient(varurl, varuser, varpass, varinsecure, vardebug)
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Creating client error %#v\n", err)
	}
	fmt.Printf("Parent ID: %s\n", varparentid)
	// }}}

	// To changing the config of debug you use this method
	// client.SwitchDebug()

	// We get all the instances that we have in our VmWare Workstation
	// AllVMs, err := client.GetAllVMs()
	// if err != nil {
	// 	log.Printf("[ERROR][MAIN] Fi: main.go Task: Listing VMs Error %#v\n", err)
	// }
	// for _, value := range AllVMs {
	// 	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nIp: %s\nProcessor: %d\nMemory: %d \n\n",
	// 		value.IdVM,
	// 		value.Path,
	// 		value.Denomination,
	// 		value.Description,
	// 		value.PowerStatus,
	// 		value.Ip,
	// 		value.CPU.Processors,
	// 		value.Memory)
	// }

	// After to read all the instances that we have in the list, we can create a new one to test it
	VM, err := client.CreateVM(varparentid, "clone-test-copy", "Test to INSERT description", 2, 1024) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Creating VMs Error %#v\n", err)
	}
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nIp: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.Ip,
		VM.CPU.Processors,
		VM.Memory)

	// Now we have one new instance, we can read which is the status
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: First Reading VM Error %#v\n", err)
	}
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nIp: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.Ip,
		VM.CPU.Processors,
		VM.Memory)

	// We want to register the instance in the UI of the VmWare Workstaion
	_, err = client.RegisterVM(VM.Denomination, VM.Path) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Registering VM Error %#v\n", err)
	}

	time.Sleep(10 * time.Second)

	// After to register the instance, we will update the values of the instance with new onece
	VM, err = client.UpdateVM(VM.IdVM, "clone-test-copy-change", "esta es una prueba de llenadao de datos", 1, 512, "off") // the id it's a test
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nIp: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.Ip,
		VM.CPU.Processors,
		VM.Memory)
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Updating VM Error %#v\n", err)
	}

	// We will read the new values
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Second Reading VM Error %#v\n", err)
	}
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nIp: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.Ip,
		VM.CPU.Processors,
		VM.Memory)

	// We want to energize the instance in order to test the PowerSwitch method
	VM, err = client.PowerSwitch(VM.IdVM, "on") // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Energizing VM Error %#v\n", err)
	}

	time.Sleep(60 * time.Second)

	// We will read the new state of the instance
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Third Reading VM Error %#v\n", err)
	}
	fmt.Printf("ID: %v\nPower Status: %s\nIP: %s \n\n",
		VM.IdVM,
		VM.PowerStatus,
		VM.Ip)

	// We shutdown the instance in order to test the PowerSwitch method
	VM, err = client.PowerSwitch(VM.IdVM, "off") // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Shutting-down VM Error %#v\n", err)
	}

	time.Sleep(60 * time.Second)

	// We want to know which was the power state of the instance
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Fourth Reading VM Error %#v\n", err)
	}
	fmt.Printf("ID: %v\nPower Status: %s \n\n",
		VM.IdVM,
		VM.PowerStatus)

	// And finally we will delete the instance that we was created
	err = client.DeleteVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("%s", err)
	}
}
