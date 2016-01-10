package faimodules

import "fmt"

func ProvisionDoc() string {

	return fmt.Sprintf("\n\nYou can provision hosts in two ways\n\n\n\tFrom the existing vlan, \n\t\tOR\n\tBy moving to the Provisioning vlan.\n\n")
}
