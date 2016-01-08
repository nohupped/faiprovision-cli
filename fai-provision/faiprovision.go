package main

import (
	"fmt"
	"os/user"
//	"os"
	"regexp"
	. "faimodules"
	"strconv"
	"strings"
)



func checkSubnet(subnet *string) (bool) {
	re := regexp.MustCompile("^(((255.){3}(255|254|252|248|240|224|192|128|0+))|((255.){2}(255|254|252|248|240|224|192|128|0+).0)|((255.)(255|254|252|248|240|224|192|128|0+)(.0+){2})|((255|254|252|248|240|224|192|128|0+)(.0+){3}))$")
	if re.MatchString(*subnet){
		return true
	} else {
		return false
	}
}

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
	fmt.Println(ProvisionDoc())
	newHost := new(Host)
	sameVlan := ValidateAndPopulateVlan(newHost)
	if sameVlan{
		localip := GetLocalIP("eth0")
		network, broadcast := GetNetworkSegment(localip)
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
				newHost.SetHostIP(IPtoCheck)
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

	}




	/*







	fmt.Println("Enter the subnet for the IP.")
	fmt.Scanln(&newHost.subnet)
	Info.Println("Subnet Entered, modified value: ", *newHost)
	if !checkSubnet(&newHost.subnet) {
		Error.Println("Subnet is not valid, dying..\n\t", *newHost)
		fmt.Println("Subnet is not valid, dying..")
		os.Exit(1)
	}

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
