package main

import (
	"fmt"
	"os/user"
//	"os"
	. "faimodules"
	"strconv"
	"strings"

)


func main() {

	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	//TODO Remove the comment below later.
/*	if currentUser.Name == "root" {
		fmt.Println("Current user is root. Run as a non-root user.")
		os.Exit(1)
	}*/
	//Opening log file
	logFile := StartLog("/var/log/fai.log", currentUser)
	defer logFile.Close()
	Info.Println("FAI Run starting as user: ", *currentUser)
	includeFiles := ReadDhcpRO("/etc/dhcp/dhcp.conf")
	Info.Println("Gathered all include files, ", includeFiles)
	fmt.Println("Gathered include files, extracting IPs from them.")

	fmt.Println(ProvisionDoc())
	newHost := new(Host)
	sameVlan := ValidateAndPopulateVlan(newHost)
	if sameVlan{
		localip := GetLocalIP("wlan0")
		network, broadcast, localmask := GetNetworkSegment(localip)
		fmt.Println(network, broadcast)
		for start := network[len(network)-1]+2; start <= broadcast[len(broadcast)-1]-1; start ++ {
			tmpnetwork := make([]int, 4, 4)
			copy(tmpnetwork, network)
			tmpnetwork[len(tmpnetwork)-1]= start
			tmpnetworkstring := make([]string, 0, 4)
			for i := 0; i<=3; i++ {
				tmp := strconv.Itoa(tmpnetwork[i])
				//fmt.Println(tmp)
				tmpnetworkstring = append(tmpnetworkstring, tmp)

			}
			IPtoCheck := strings.Join(tmpnetworkstring[:], ".")
			Info.Println("Checking if ", IPtoCheck, "is alive")
			alive := DoPing(IPtoCheck)
			if alive != "" {
				fmt.Println(alive)
				Warning.Println(alive, "checking nextIP")
				fmt.Println("Checking next")
			}else {
				Info.Println(IPtoCheck, "is not alive in ping check.")
				//TODO Add a check to verify ip not in dhcp.conf
				newHost.SetHostIP(IPtoCheck)
				newHost.SetHostSubnetInt(localmask)
				Info.Println("Host IP to be configured as", IPtoCheck)
				break
			}

		}
		fmt.Println("IP set as ", newHost.GetIP())
		ValidateAndPopulateHostname(newHost)
		ValidateAndPopulateMacID(newHost)


	}else {
		fmt.Println("todo for manual IP provision")
		ValidateAndPopulateIP(newHost)
		ValidateAndPopulateSubnet(newHost)
	}




	/*









	for ; ;  {
		fmt.Println("type 'yes' for recloning an existing production machine, or 'no' for installing to a fresh server")
		var tmp string
		fmt.Scanln(&tmp)
		Info.Println("type 'no' for recloning an existing production machine, or 'yes' for installing to a fresh server: " , tmp)
		if tmp == "yes" || tmp == "no"{
			switch tmp {
			case "yes":
				&newHost.reClone=true
			case "no":
				&newHost.reClone=false
			}
			break
		}else {
			Error.Println("Typed value is '", tmp, "', type 'yes' or 'no'. ")
		}
	}*/



}
