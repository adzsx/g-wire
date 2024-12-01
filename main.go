package main

import (
	"github.com/adzsx/gwire/cmd"
)

func main() {
	cmd.Execute()
}

// func main() {

// 	log.SetFlags(0)
// 	args := os.Args

// 	input := utils.Format(args)

// 	if input.Action == "version" {
// 		log.Println(version)
// 		os.Exit(0)
// 	}

// 	err := utils.CheckInput(input)
// 	utils.Err(err, true)

// 	if len(args) < 3 && input.Action != "help" && input.Action != "info" {
// 		log.Println("Enter --help for help")
// 		os.Exit(0)
// 	} else if input.Action == "help" {
// 		log.Print(help)
// 		os.Exit(0)
// 	}

// 	if input.Action == "listen" {
// 		netcli.HostSetup(input)

// 	} else if input.Action == "connect" {

// 		netcli.ClientSetup(input)
// 	} else if input.Action == "info" {

// 		ip, mask, nHosts, _ := netcli.Info()
// 		log.Printf("Private IP: 		%v\nSubnetmask: 		%v\nNumber of hosts: 	%v", ip, mask, nHosts)

// 	}

// }
