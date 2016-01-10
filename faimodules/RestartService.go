package faimodules

import (
	"fmt"
	"os/exec"
	"bytes"
)

func RestartDHCP(initscript string) error  {
	fmt.Println("This program will attempt to restart the service")
	Info.Println("This program will attempt to restart the service")
	cmd := exec.Command("sudo", initscript, "restart")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err, stderr.String())
		Error.Println(err, stderr.String())
		return err
	}
	fmt.Println(out.String())
	Info.Println(out.String())
	return err
}

