package main

import (
	"log"
	"os"

	"github.com/adzsx/gwire/pkg/netcli"
	"github.com/adzsx/gwire/pkg/utils"
)

var (
	help = `
gwire usage:
	gwire [mode] [flags]

Modes:
	[No input]	Uses mode to connect and chat
	info		Shows network related info


Flags:
	    --help			help message
	-l, --listen			listen
	-p, --port 	[port]		use port [port]
	-h, --host 	[host]		Connect to [host]-(Ip)
	-v, --verbose			Show some more info
	-u, --username 	[username]	[username] is didsplayed for other users
	-t, --time			enable timestamps
	-s, --slowmode	[seconds]	Enable slowmode
	-e, --encrypt	([password])	If [password] is given, use AES, if not, encrypt automatic with RSA
	`

	version = "gwire v1.3.0"
)

func main() {

	log.SetFlags(0)
	args := os.Args

	input := utils.Format(args)

	if input.Action == "version" {
		log.Println(version)
		os.Exit(0)
	}

	if len(args) < 3 && input.Action != "help" && input.Action != "info" {
		log.Println("Enter --help for help")
		os.Exit(0)
	} else if input.Action == "help" {
		log.Print(help)
		os.Exit(0)
	}

	if input.Action == "listen" {
		netcli.HostSetup(input)

	} else if input.Action == "connect" {

		netcli.ClientSetup(input)
	} else if input.Action == "info" {

		ip, mask, nHosts, _ := netcli.Info()
		log.Printf("Private IP: 		%v\nSubnetmask: 		%v\nNumber of hosts: 	%v", ip, mask, nHosts)

	}

}
