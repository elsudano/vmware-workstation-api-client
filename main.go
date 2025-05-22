package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/elsudano/vmware-workstation-api-client/wsapiclient"
)

var paragraph_color = color.Blue
var title_line_color = color.White
var value_color = color.Yellow

func PrintVM(VM *wsapiclient.MyVm) {
	fmt.Println(
		color.Ize(title_line_color, " ID:"), color.Ize(value_color, VM.IdVM), "\n",
		color.Ize(title_line_color, "Path:"), color.Ize(value_color, VM.Path), "\n",
		color.Ize(title_line_color, "Denomination:"), color.Ize(value_color, VM.Denomination), "\n",
		color.Ize(title_line_color, "Description:"), color.Ize(value_color, VM.Description), "\n",
		color.Ize(title_line_color, "Processor:"), color.Ize(value_color, VM.CPU.Processors), "\n",
		color.Ize(title_line_color, "Memory:"), color.Ize(value_color, VM.Memory), "\n",
		color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus), "\n",
		color.Ize(title_line_color, "HostName:"), color.Ize(value_color, VM.DNS.Hostname), "\n",
		color.Ize(title_line_color, "DomainName:"), color.Ize(value_color, VM.DNS.Domainname), "\n",
		color.Ize(title_line_color, "DNServers:"), color.Ize(value_color, VM.DNS.Servers),
	)
	if VM.NICS != nil {
		for nic := range VM.NICS {
			fmt.Println(
				color.Ize(title_line_color, " Mac:"), color.Ize(value_color, VM.NICS[nic].Mac), "\n",
				color.Ize(title_line_color, "Ip:"), color.Ize(value_color, VM.NICS[nic].Ip),
			)
		}
	}
}

func main() {
	file, err := os.Open("config.ini")
	if err != nil {
		log.Fatalf("Failed opening file %s, please make sure the config file exists", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var varuser, varpass, varurl, varparentid, vardebug string
	var varinsecure bool
	for scanner.Scan() {
		array := strings.SplitN(scanner.Text(), ":", 2)
		key := strings.ToLower(array[0])
		value := array[1]
		if !strings.HasPrefix(key, "#") {
			fmt.Println(
				color.Ize(title_line_color, "Temp: "), color.Ize(value_color, array),
				color.Ize(title_line_color, "Key: "), color.Ize(value_color, key),
				color.Ize(title_line_color, "Value: "), color.Ize(value_color, value),
			)
			switch key {
			case "user":
				varuser = strings.TrimSpace(value)
			case "password":
				varpass = strings.TrimSpace(value)
			case "baseurl":
				varurl = strings.TrimSpace(value)
			case "parentid":
				varparentid = strings.TrimSpace(value)
			case "insecure":
				if strings.TrimSpace(value) == "true" {
					varinsecure = true
				} else {
					varinsecure = false
				}
			case "debug":
				vardebug = strings.TrimSpace(value)
			}
		}
	}
	fmt.Println(color.Ize(paragraph_color, "First of all we want to see which are the Input Values:"))
	fmt.Println(
		color.Ize(title_line_color, "varuser:"), color.Ize(value_color, varuser),
		color.Ize(title_line_color, "varpass:"), color.Ize(value_color, varpass),
		color.Ize(title_line_color, "varurl:"), color.Ize(value_color, varurl),
		color.Ize(title_line_color, "varparentid:"), color.Ize(value_color, varparentid),
		color.Ize(title_line_color, "varinsecure:"), color.Ize(value_color, varinsecure),
		color.Ize(title_line_color, "vardebug:"), color.Ize(value_color, vardebug),
	)
	fmt.Println()
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
		os.Exit(10)
	}
	fmt.Println(color.Ize(paragraph_color, "We can see here the value of the ParentID VM:"))
	fmt.Println(color.Ize(title_line_color, "Parent ID:"), color.Ize(value_color, varparentid))
	// }}}

	// To changing the config of debug you use this method
	// client.SwitchDebugLevel("NONE")

	// We get all the instances that we have in our VmWare Workstation
	fmt.Println(color.Ize(paragraph_color, "Now we going to list all the VMs that we have in the VMWare Workstation:"))
	AllVMs, err := client.GetAllVMs()
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Listing VMs Error %#v\n", err)
		// os.Exit(11)
	}
	for _, item := range AllVMs {
		if item.IdVM != "II82IG14IV0UHB7QHAI3J0G44RJBGNR5" { // is just because this VM is encripted
			VM, err := client.LoadVM(item.IdVM)
			if err != nil {
				log.Printf("[ERROR][MAIN] Fi: main.go Task: First time Reading VM Error %#v\n", err)
				os.Exit(13)
			}
			PrintVM(VM)
		}
	}
	// After to read all the instances that we have in the list, we can create a new one to test it
	fmt.Println(color.Ize(paragraph_color, "We are using the CreateVM method to create the VM test."))
	VM, err := client.CreateVM(varparentid, "clone-test-copy", "Test to INSERT description", 2, 1024)
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Creating VMs Error %#v\n", err)
		os.Exit(12)
	}
	// Now we have one new instance, we can read which is the status
	fmt.Println(color.Ize(paragraph_color, "Now we going to review the state of the VM that we have created:"))
	VM, err = client.LoadVM(VM.IdVM)
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: First time Reading VM Error %#v\n", err)
		os.Exit(13)
	}
	PrintVM(VM)
	// We want to register the instance in the UI of the VmWare Workstaion
	// fmt.Println(color.Ize(paragraph_color, "After that, we have registered in the UI of VMWare Workstation"))
	// _, err = client.RegisterVM(VM.Denomination, VM.Path)
	// if err != nil {
	// 	log.Printf("[ERROR][MAIN] Fi: main.go Task: Registering VM Error %#v\n", err)
	// 	os.Exit(14)
	// }
	// time.Sleep(10 * time.Second)
	// After to register the instance, we will update the values of the instance with new onece
	fmt.Println(color.Ize(paragraph_color, "Now, we going to Energize the VM with the UpdateVM method."))
	VM, err = client.UpdateVM(VM.IdVM, "clone-test-copy-change", "esta es una prueba de llenadao de datos", 2, 1024, "on")
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Updating VM Error %#v\n", err)
		os.Exit(15)
	}
	time.Sleep(60 * time.Second) // we need to wait because the VM take time to be ready
	// fmt.Println(color.Ize(paragraph_color, "We can confirm that the VM is working properly:"))
	// VM, err = client.LoadVM(VM.IdVM)
	// if err != nil {
	// 	log.Printf("[ERROR][MAIN] Fi: main.go Task: Second time Reading VM Error %#v\n", err)
	// 	os.Exit(16)
	// }
	PrintVM(VM)
	// We want to shutdown the instance in order to test the UpdateVM method
	fmt.Println(color.Ize(paragraph_color, "Now, we going to shutdown and change the propierties of the VM with the UpdateVM method"))
	VM, err = client.UpdateVM(VM.IdVM, "clone-test-copy-change", "esta es una prueba de llenadao de datos", 1, 512, "off")
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Updating VM Error %#v\n", err)
		os.Exit(17)
	}
	time.Sleep(30 * time.Second)
	// fmt.Println(color.Ize(paragraph_color, "We confirm that the VM is off and it's propierties has changed:"))
	// VM, err = client.LoadVM(VM.IdVM)
	// if err != nil {
	// 	log.Printf("[ERROR][MAIN] Fi: main.go Task: Third time Reading VM Error %#v\n", err)
	// 	os.Exit(18)
	// }
	PrintVM(VM)
	// And finally we will delete the instance that we was created
	fmt.Println(color.Ize(paragraph_color, "And finally we have deleted the VM"))
	err = client.DeleteVM(VM.IdVM)
	if err != nil {
		log.Printf("%s", err)
		os.Exit(19)
	}
	fmt.Println()
}
