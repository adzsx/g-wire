package client

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/adzsx/gwire/internal/netutils"
	"github.com/adzsx/gwire/internal/utils"
)

// Function connects to host with TCP
func ClientSetup(host string, port string, handshake bool, enc string, username string, time bool) {
	log.SetFlags(0)
	if time {
		log.SetFlags(log.Ltime)
	}

	inpHost = host
	inpPort = port
	inpHandshake = handshake
	inpEnc = enc
	inpUsername = username
	displayTime = time

	// Connect to host

	if host != "scan" {
		conn, err = net.Dial("tcp", host+":"+port)
	} else {
		// Scan every host in network for open port
		hosts, err := netutils.GetHosts(netutils.Subnet())

		utils.Err(err, true)
		host, conn = netutils.ScanRange(hosts, port)
	}

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		utils.Err(errors.New("connection refused by destination"), true)
	}

	utils.Print("Connected to "+host+":"+port+"\n", 0)

	if handshake {
		if enc == "auto" {
			err = initClient(conn)
			utils.Err(err, true)
		}

		utils.Print("Setup finished", 1)
		client(conn)
	} else {

	}
}
