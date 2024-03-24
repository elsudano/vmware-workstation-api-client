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

func main() {
	paragraph_color := color.Blue
	title_line_color := color.White
	value_color := color.Bold
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
		fmt.Println(
			color.Ize(title_line_color, "Temp: "), color.Ize(value_color, temp),
			color.Ize(title_line_color, "Key: "), color.Ize(value_color, vtemp),
			color.Ize(title_line_color, "Value: "), color.Ize(value_color, vtemp2),
		)
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
	}
	fmt.Println(color.Ize(paragraph_color, "We can see here the value of the ParentID VM:"))
	fmt.Println(color.Ize(title_line_color, "Parent ID:"), color.Ize(value_color, varparentid))
	// }}}

	// To changing the config of debug you use this method
	// client.SwitchDebug()

	// We get all the instances that we have in our VmWare Workstation
	AllVMs, err := client.GetAllVMs()
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Listing VMs Error %#v\n", err)
	}
	fmt.Println(color.Ize(paragraph_color, "Now we go to list all teh VMs that we have in the VMWare Workstation:"))
	for _, VM := range AllVMs {
		fmt.Println(color.Ize(title_line_color, "ID:"), color.Ize(value_color, VM.IdVM), "\n",
			color.Ize(title_line_color, "Path:"), color.Ize(value_color, VM.Path), "\n",
			color.Ize(title_line_color, "Denomination:"), color.Ize(value_color, VM.Denomination), "\n",
			color.Ize(title_line_color, "Description:"), color.Ize(value_color, VM.Description), "\n",
			color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus), "\n",
			color.Ize(title_line_color, "Ip:"), color.Ize(value_color, VM.Ip), "\n",
			color.Ize(title_line_color, "Processor:"), color.Ize(value_color, VM.CPU.Processors), "\n",
			color.Ize(title_line_color, "Memory:"), color.Ize(value_color, VM.Memory),
		)
	}

	// After to read all the instances that we have in the list, we can create a new one to test it
	VM, err := client.CreateVM(varparentid, "clone-test-copy", "Test to INSERT description", 1, 512) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Creating VMs Error %#v\n", err)
	}
	fmt.Println(color.Ize(paragraph_color, "We going to create the first VM in VMWare Workstation:"))
	fmt.Println(color.Ize(title_line_color, "ID:"), color.Ize(value_color, VM.IdVM), "\n",
		color.Ize(title_line_color, "Path:"), color.Ize(value_color, VM.Path), "\n",
		color.Ize(title_line_color, "Denomination:"), color.Ize(value_color, VM.Denomination), "\n",
		color.Ize(title_line_color, "Description:"), color.Ize(value_color, VM.Description), "\n",
		color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus), "\n",
		color.Ize(title_line_color, "Ip:"), color.Ize(value_color, VM.Ip), "\n",
		color.Ize(title_line_color, "Processor:"), color.Ize(value_color, VM.CPU.Processors), "\n",
		color.Ize(title_line_color, "Memory:"), color.Ize(value_color, VM.Memory),
	)
	fmt.Println()
	// Now we have one new instance, we can read which is the status
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: First Reading VM Error %#v\n", err)
	}
	fmt.Println(color.Ize(paragraph_color, "Now we go to review the state of the VM that we have created:"))
	fmt.Println(color.Ize(title_line_color, "ID:"), color.Ize(value_color, VM.IdVM), "\n",
		color.Ize(title_line_color, "Path:"), color.Ize(value_color, VM.Path), "\n",
		color.Ize(title_line_color, "Denomination:"), color.Ize(value_color, VM.Denomination), "\n",
		color.Ize(title_line_color, "Description:"), color.Ize(value_color, VM.Description), "\n",
		color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus), "\n",
		color.Ize(title_line_color, "Ip:"), color.Ize(value_color, VM.Ip), "\n",
		color.Ize(title_line_color, "Processor:"), color.Ize(value_color, VM.CPU.Processors), "\n",
		color.Ize(title_line_color, "Memory:"), color.Ize(value_color, VM.Memory),
	)
	fmt.Println()

	// We want to register the instance in the UI of the VmWare Workstaion
	fmt.Println(color.Ize(paragraph_color, "After that, we have registered in the UI of VMWare Workstation"))
	_, err = client.RegisterVM(VM.Denomination, VM.Path) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Registering VM Error %#v\n", err)
	}
	fmt.Println()

	time.Sleep(10 * time.Second)

	// After to register the instance, we will update the values of the instance with new onece
	fmt.Println(color.Ize(paragraph_color, "Now, we going to Powerize the VM with the UpdateVM method:"))
	VM, err = client.UpdateVM(VM.IdVM, "clone-test-copy-change", "esta es una prueba de llenadao de datos", 2, 1024, "on") // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Updating VM Error %#v\n", err)
	}
	fmt.Println()

	time.Sleep(60 * time.Second)

	// We will read the new values
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Second Reading VM Error %#v\n", err)
	}
	fmt.Println(color.Ize(paragraph_color, "We can confirm that the VM is working properly:"))
	fmt.Println(color.Ize(title_line_color, "ID:"), color.Ize(value_color, VM.IdVM), "\n",
		color.Ize(title_line_color, "Path:"), color.Ize(value_color, VM.Path), "\n",
		color.Ize(title_line_color, "Denomination:"), color.Ize(value_color, VM.Denomination), "\n",
		color.Ize(title_line_color, "Description:"), color.Ize(value_color, VM.Description), "\n",
		color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus), "\n",
		color.Ize(title_line_color, "Ip:"), color.Ize(value_color, VM.Ip), "\n",
		color.Ize(title_line_color, "Processor:"), color.Ize(value_color, VM.CPU.Processors), "\n",
		color.Ize(title_line_color, "Memory:"), color.Ize(value_color, VM.Memory),
	)
	fmt.Println()

	// We want to energize the instance in order to test the PowerSwitch method
	fmt.Println(color.Ize(paragraph_color, "Now, we going to shutdown the VM with the PowerSwich method"))
	VM, err = client.PowerSwitch(VM.IdVM, "off") // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Energizing VM Error %#v\n", err)
	}
	fmt.Println()

	time.Sleep(60 * time.Second)

	// We will read the new state of the instance
	VM, err = client.ReadVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("[ERROR][MAIN] Fi: main.go Task: Third Reading VM Error %#v\n", err)
	}
	fmt.Println(color.Ize(paragraph_color, "We confirm that the VM is off:"))
	fmt.Println(color.Ize(title_line_color, "ID:"), color.Ize(value_color, VM.IdVM), "\n",
		color.Ize(title_line_color, "Power Status:"), color.Ize(value_color, VM.PowerStatus),
	)
	fmt.Println()

	// And finally we will delete the instance that we was created
	fmt.Println(color.Ize(paragraph_color, "And finally we destroy the VM"))
	err = client.DeleteVM(VM.IdVM) // the id it's a test
	if err != nil {
		log.Printf("%s", err)
	}
	fmt.Println()
}
