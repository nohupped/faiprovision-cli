
package faimodules

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

func DoPing(ip string) string {
	var ipAlive string

	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		ipAlive = fmt.Sprintf("IP Addr: %s Alive, RTT: %v\n", addr.String(), rtt)
	}
	//p.OnIdle = func() {
	//	fmt.Println("finished")
	//}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
		Error.Println(err)
		os.Exit(1)
	}
	return ipAlive
}