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
	var vardebug bool
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
		if vtemp == "debug" && strings.TrimSpace(temp[1]) == "true" {
			vardebug = true
		} else {
			vardebug = false
		}
	}
	file.Close()

	client, err := wsapiclient.New()
	client.ConfigCli(varurl, varuser, varpass, vardebug)
	AllVMs, err := client.GetAllVMs()
	if err != nil {
		log.Fatalf("%s", err)
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
	VM, err := client.GetVM(AllVMs[0].IdVM)
	fmt.Printf("VM Id: %v\nVM Processors: %d\nVM Memory: %d\nVM Status: %v\n",
		VM.IdVM, VM.CPU.Processors, VM.Memory, VM.PowerStatus)
	if err != nil {
		log.Fatalf("%s", err)
	}
}
