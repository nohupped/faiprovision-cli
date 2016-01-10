package faimodules

import (
	"strconv"
	"strings"
)

type Host struct  {
	hostname string
	macID string
	ip string
	subnet string
	gateway string
	sameVlan bool
}

func (h Host) GetHostname() string {
	return h.hostname
}

func (h Host) GetMacID() string {
	return h.macID
}

func (h Host) GetIP() string  {
	return h.ip
}

func (h Host) GetSubnet() string  {
	return h.subnet
}

func (h Host) IfSameVlan() bool  {
	return h.sameVlan
}

func (h *Host) SetHostname(x string)  {
	h.hostname = x
}
func (h *Host) SetSameVlan(x bool)  {
	h.sameVlan = x
}
func (h *Host) SetHostIP(x string){
	h.ip = x
}
func (h *Host) SetHostSubnetInt(x []int)  {
	s := make([]string, 0, 4)
	for _, i := range x{
		s = append(s, strconv.Itoa(i))
	}
	h.subnet = strings.Join(s[:], ".")
}

//SetHostDefaultRouteInt will accept network as the parameter, adds 1 to the host part, and assigns it to h.gateway
func (h *Host) SetHostDefaultRouteInt(x []int)  {
	route_tmp := x[len(x)-1]+1
	//route := make([]int 1, 1)
	route := make([]int, 4)
	copy(route, x)
	route[len(route)-1] = route_tmp
	route_string := make([]string, 0, 4)
	for _, i := range route  {
		tmp := strconv.Itoa(i)
		route_string = append(route_string, tmp)
	}
	h.gateway = strings.Join(route_string[:], ".")
}
func (h *Host) GetRoute() string  {
	return h.gateway
}

