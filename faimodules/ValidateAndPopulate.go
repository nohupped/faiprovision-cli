package faimodules

import (
	"fmt"
	"regexp"
	"os"
)

func ValidateAndPopulateVlan(h *Host) bool {
	Info.Println("Initiated Struct Host with default values", *h)
	fmt.Println("Welcome to interactive FAI provisioning. You can track the progress from DRAC.")
	for ; ;  {
		fmt.Println("Type 'yes' if the new host is on the same vlan as that of the provision server, 'no' if otherwise")
		var tmp string
		fmt.Scanln(&tmp)
		Info.Println("Type 'yes' if the new host is on the same vlan as that of the provision server, 'no' if otherwise: " , tmp)
		if tmp == "yes" || tmp == "no"{
			switch tmp {
			case "yes":
				h.SetSameVlan(true)
				Info.Println("Host to be provisioned is set to be on the same vlan as this server.")
				return true
			case "no":
				h.SetSameVlan(false)
				Info.Println("Host to be provisioned is set to be on a different vlan as this server.")
				return false
			}
		//	break
		}else {
			Error.Println("Typed value is '", tmp, "', type 'yes' or 'no'. ")
		}
	}

}

func checkHostName(hostname *string) (bool){
	re := regexp.MustCompile("^[a-zA-Z0-9_.]*$")
	if re.MatchString(*hostname) {
		return true
	} else {
		return false
	}
}

func ValidateAndPopulateHostname(h *Host)  {
	fmt.Println("What is your new hostname? If unsure now, give IP.")
	fmt.Scanln(&h.hostname)
	Info.Println("Hostname value entered, modified value: ", *h)

	if !checkHostName(&h.hostname){
		Error.Println("Hostname doesn't comply with standards, Dying.. \n\t", *h)
		fmt.Println("Hostname doesn't comply with standards, dying...")
		os.Exit(1)
	}
}


func checkMacid(macid *string) (bool) {
	re := regexp.MustCompile("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$")
	if re.MatchString(*macid){
		return true
	}else {
		return false
	}
}

func ValidateAndPopulateMacID(h *Host)  {
	fmt.Println("Enter the macId of the new host that will be PXE'd.")
	fmt.Scanln(&h.macID)
	Info.Println("MacID Entered, modified value: ", *h)
	if !checkMacid(&h.macID) {
		Error.Println("MacID doesn't comply with standards, dying.. \n\t", *h)
		fmt.Println("MacID doesn't comply with standards, dying...")
		os.Exit(1)
	}
}


func checkIP(ip *string) (bool) {
	re := regexp.MustCompile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]).){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
	if re.MatchString(*ip) {
		return true
	} else {
		return false
	}
}

func ValidateAndPopulateIP(h *Host)  {
	fmt.Println("Enter the IP address to be assigned to the host.")
	fmt.Scanln(&h.ip)
	Info.Println("IP address Entered, modified value: ", *h)
	if !checkIP(&h.ip) {
		Error.Println("IP is not valid, dying..\n\t", *h)
		fmt.Println("IP is not valid, dying..")
		os.Exit(1)
	}
}


func checkSubnet(subnet *string) (bool) {
	re := regexp.MustCompile("^(((255.){3}(255|254|252|248|240|224|192|128|0+))|((255.){2}(255|254|252|248|240|224|192|128|0+).0)|((255.)(255|254|252|248|240|224|192|128|0+)(.0+){2})|((255|254|252|248|240|224|192|128|0+)(.0+){3}))$")
	if re.MatchString(*subnet){
		return true
	} else {
		return false
	}
}

func ValidateAndPopulateSubnet(h *Host)  {
	fmt.Println("Enter the subnet for the IP.")
	fmt.Scanln(&h.subnet)
	Info.Println("Subnet Entered, modified value: ", *h)
	if !checkSubnet(&h.subnet) {
		Error.Println("Subnet is not valid, dying..\n\t", *h)
		fmt.Println("Subnet is not valid, dying..")
		os.Exit(1)
	}
}