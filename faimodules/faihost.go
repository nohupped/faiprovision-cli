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
	reClone bool
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

func (h Host) GetRecloneStats() bool  {
	return h.reClone
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

