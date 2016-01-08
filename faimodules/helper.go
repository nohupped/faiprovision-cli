package faimodules

import "fmt"

func ProvisionDoc() string {

	return fmt.Sprintf("You can provision hosts in two ways\n\tFrom the existing vlan, \n\t\tOR\n\tBy moving to the Provisioning vlan.\n")
}
