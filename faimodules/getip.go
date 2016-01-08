package faimodules

import (
	"fmt"
	"net"
	"strings"
	"strconv"
)

func GetLocalIP(nic string) string{

	fmt.Println("Getting Local IP")
	Info.Println("Since the host to be provisioned is set to be on the same vlan, auto-assigning IP.")
	var localip []net.Addr
	ifaces, err := net.Interfaces()

	if err != nil {
		panic(err)
	}
	Info.Println("Gathered local interface details to determine network and subnet", ifaces)


	for _, iface := range ifaces{
		if iface.Name == nic{
			localip, _ = iface.Addrs()
		}
	}
	fmt.Println(localip[0].String())
	Info.Println("Choosing ", localip[0].String(), "according to the parameter")

	return localip[0].String()

}

func GetNetworkSegment(localIPADDR string) (NetworkAddr []int, BroadcastAddr []int)  {
	fmt.Println("Getting  network and broadcast address of ", localIPADDR)
	Info.Println("Stripping down IP to network")
	ipsubnetSlice := strings.Split(localIPADDR, "/")
	ipSlice := strings.Split(ipsubnetSlice[0], ".")
	localip := make([]int, 0, 4)
	for _, ip := range ipSlice{
		tmp, err := strconv.Atoi(ip)
		if err != nil {
			Error.Println(err)
			panic(err)
		}
		localip = append(localip, tmp)
	}


	netmaskint, err := strconv.ParseUint(ipsubnetSlice[1], 10, 64)

	if err != nil {
	panic(err)
	}
	var mask uint64
	mask = (0xFFFFFFFF << (32 - netmaskint)) & 0xFFFFFFFF;
	localmask := make([]int, 0, 4)
	var dmask uint64
	dmask = 32
	for i := 1; i <= 4; i++{
		Info.Println("Bitwise shifting of ", mask , "and logical AND ", 0xFF, "to convert CIDR to netmask")
		tmp := mask >> (dmask - 8) & 0xFF
		Info.Println(tmp)
		tmp1 := int(tmp)
		localmask = append(localmask, tmp1)
		dmask -= 8
	}
	//fmt.Printf("%lu.%lu.%lu.%lu\n", mask >> 24, (mask >> 16) & 0xFF, (mask >> 8) & 0xFF, mask & 0xFF);
	Info.Println("Converted to 32 bit mask", localmask)

	Network := make([]int, 0, 4)
	for i := 0; i <= 3; i++ {
		Network = append(Network, localip[i] & localmask[i])
	}
	Info.Println("Calculated network for the current network ", localIPADDR, " as ",  Network)
	Info.Println("Calculated netmask for the current network ",Network, "as",  localmask)
	Broadcast := make([]int, 0, 0)

	invertedmask := make([]int, 0, 4)
	for i := 0; i <= 3; i++{
		var tmp uint8
		tmp = ^uint8(localmask[i])
		tmp1 := int(tmp)
		invertedmask = append(invertedmask, tmp1)
	}
	Info.Println("Flipping bits", invertedmask)

	for i := 0; i <= 3; i++{
		Broadcast = append(Broadcast, Network[i] | invertedmask[i])
		Info.Println("Bitwise OR between ", Network[i], " and ", invertedmask[i], ":", Network[i] | invertedmask[i])
	}
	Info.Println("Calculated broadcast IP for the network",Broadcast)

	return Network, Broadcast
}