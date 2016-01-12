package main

import (
	"fmt"
	"os/user"
	. "faimodules"
	"strconv"
	"strings"
	"os"
	"syscall"
)


func main() {
	currentUser, err := user.Current()
	CheckForError(err)

	if currentUser.Name == "root" {
		fmt.Println("Current user is root. Run as a non-root user.")
		os.Exit(1)
	}

	//Opening log file
	logFile := StartLog("/var/log/fai.log", currentUser)
	defer logFile.Close()
	includepath := "/etc/dhcp/hosts/"
	dhcpmainconf := "/etc/dhcp/dhcpd.conf"
	nextserverIP := "172.20.17.106"
	DHCPInitScript := "/etc/init.d/isc-dhcp-server"
	ProgramLock := "/tmp/.fai.lock"
	networkinterface := "eth0"

	lockProgram, err := os.Create(ProgramLock)
	defer os.Remove(ProgramLock)
	CheckForError(err)
	defer lockProgram.Close()
	GetLock(lockProgram, syscall.LOCK_EX)
	defer UngetLock(lockProgram)

	Info.Println("FAI Run starting as user: ", *currentUser)
	includeFiles := ReadDhcpRO(dhcpmainconf)
	Info.Println("Gathered all include files, ", includeFiles)
	fmt.Println("Gathered include files, extracting IPs from them for checking later.")
	dhcpips := GetIpFromInclude(includeFiles)

	fmt.Println(ProvisionDoc())
	newHost := new(Host)
	sameVlan := ValidateAndPopulateVlan(newHost)
	if sameVlan{
		localip := GetLocalIP(networkinterface)
		network, broadcast, localmask := GetNetworkSegment(localip)
		fmt.Println(network, broadcast)
		Info.Println("Setting default gateway")
		newHost.SetHostDefaultRouteInt(network)
		Info.Println("Default gateway set, modified value: ", newHost)
		fmt.Println("Default gw set as, ", newHost.GetRoute())
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
			}else if InSlice(IPtoCheck, dhcpips){
				fmt.Println(IPtoCheck, "found in dhcp conf")
				Warning.Println(IPtoCheck, "found in dhcp conf, checking next")
			} else {
				Info.Println(IPtoCheck, "is not alive in ping check.")
				newHost.SetHostIP(IPtoCheck)
				newHost.SetHostSubnetInt(localmask)
				Info.Println("Host IP to be configured as", IPtoCheck)
				Info.Println("Host Netmask to be configured as", localmask)
				break
			}

		}
		fmt.Println("IP set as ", newHost.GetIP())
		ValidateAndPopulateHostname(newHost)
		ValidateAndPopulateMacID(newHost)

	}else {
		fmt.Println("Manual configuration")
		ValidateAndPopulateHostname(newHost)
		ValidateAndPopulateIP(newHost)
		ValidateAndPopulateSubnet(newHost)
		ValidateAndPopulateMacID(newHost)
		ValidateAndPopulateRoute(newHost)

	}



	FinalSet := fmt.Sprintf("%+v\n",*newHost)
	Info.Println("The configuration about to push is", FinalSet)
	fmt.Println("The final configuration is, ", FinalSet, "Enter 'yes' to continue, any other key to abort")
	var tmp string
	fmt.Scanln(&tmp)
	if tmp != "yes"{
		Info.Println("Aborting upon user request.")
		os.Exit(1)
	}
	fmt.Println("Taking Backup")
	Backup := TakeBackup(dhcpmainconf)
	fmt.Println(Backup)
	includeconf := includepath + newHost.GetHostname() + ".conf"
	fmt.Println(includeFiles)
	for _, check := range includeFiles{
		Info.Println("Comparing configfiles", check, includeconf)
		if check == includeconf{
			fmt.Println(includeconf, "already exists, please use a different hostname.")
			Error.Println(includeconf, "already exists, please use a different hostname.")
			os.Exit(1)
		}
	}
	Info.Println("Opening DHCP configuration to add include file entry")
	WriteIncludeToMainConf(includeconf, dhcpmainconf)
	WriteToIncludeConf(includeconf, newHost, nextserverIP)
	err = RestartDHCP(DHCPInitScript)
	if err != nil {
		CopyFiles(Backup, dhcpmainconf)
		os.Remove(includeconf)
		fmt.Println("Cleanup and restore done\n\nAttempting to restart service again. If the service fail to start again, please escalate.")
		err = RestartDHCP(DHCPInitScript)
		CheckForError(err)
	} else {
		Info.Println("Settings added to DHCP and restarted.")
		fmt.Println("Settings successfully added to dhcp server and restarted. Please pxe boot the new server.")
	}

}
