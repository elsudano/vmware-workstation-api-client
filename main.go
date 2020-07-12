package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
)

func main() {
	file, err := os.Open("config.ini")
	if err != nil {
		log.Fatalf("Failed opening file %s, please make sure the config file exists", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var varuser, varpass, varurl string
	var varhttp, vardebug bool
	for scanner.Scan() {
		temp := strings.SplitN(scanner.Text(), ":", 2)
		vtemp := strings.ToLower(temp[0])
		if vtemp == "user" {
			varuser = strings.TrimSpace(temp[1])
		}
		if vtemp == "password" {
			varpass = strings.TrimSpace(temp[1])
		}
		if vtemp == "baseurl" {
			varurl = strings.TrimSpace(temp[1])
		}
		if vtemp == "https" && strings.TrimSpace(temp[1]) == "true" {
			varhttp = true
		} else {
			varhttp = false
		}
		if vtemp == "debug" && strings.TrimSpace(temp[1]) == "true" {
			vardebug = true
		} else {
			vardebug = false
		}
	}
	file.Close()

	// First option we will can create a client with the default values {{{
	client, err := wsapiclient.New()
	client.ConfigCli(varurl, varuser, varpass, varhttp, vardebug)
	// }}}

	// Second option we will setting the values in the creation moment {{{
	// client, err := wsapiclient.NewClient(varurl, varuser, varpass, varhttp, vardebug)
	// }}}

	// To changing the config of debug you use this method
	client.SwitchDebug()

	// AllVMs, err := client.GetAllVMs()
	// if err != nil {
	// 	log.Fatalf("%s", err)
	// }
	// for _, value := range AllVMs {
	// 	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nProcessor: %d\nMemory: %d \n\n",
	// 		value.IdVM,
	// 		value.Path,
	// 		value.Denomination,
	// 		value.Description,
	// 		value.PowerStatus,
	// 		value.CPU.Processors,
	// 		value.Memory)
	// }

	VM, err := client.CreateVM("PQLRP873B1P5FL1KFL5P3NN64CS5M9A3", "clone-test-copy") // the id it's a test
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.CPU.Processors,
		VM.Memory)
	if err != nil {
		log.Fatalf("%s", err)
	}

	VM, err = client.ReadVM("0D6IIC428UP5DG1631T5N4E5VTEJRUVD") // the id it's a test
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.CPU.Processors,
		VM.Memory)
	if err != nil {
		log.Fatalf("%s", err)
	}

	VM, err = client.UpdateVM(VM.IdVM, "clone-test-copy-change", "esta es una prueba de llenadao de datos", 2, 512) // the id it's a test
	fmt.Printf("ID: %v\nPath: %v\nDenomination: %v\nDescription: %s\nPower Status: %s\nProcessor: %d\nMemory: %d \n\n",
		VM.IdVM,
		VM.Path,
		VM.Denomination,
		VM.Description,
		VM.PowerStatus,
		VM.CPU.Processors,
		VM.Memory)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = client.DeleteVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Fatalf("%s", err)
	}
}
