package host

import (
	"log"
	"os"

	"github.com/adzsx/gwire/internal/utils"
	"github.com/adzsx/gwire/pkg/crypt"
)

// Set up listener for each port on list
func HostSetup(inpPort string, inpHandshake bool, inpEnc string, inpUsername string, time bool) {
	port = inpPort
	handshake = inpHandshake
	enc = inpEnc
	username = inpUsername
	displayTime = time
	// Global slice for distributing messages
	var message = [][]string{}

	if enc == "auto" {
		auto = true
		var err error
		utils.Print("Generating password\n", 2)
		enc, err = crypt.GenPasswd()
		utils.Err(err, true)
	}
	// Set up listener for every port in range
	for _, port := range port {

		// wg = WaitGroup (Variable to wait until variable hits 0)
		wg.Add(1)

		go connSetup(string(port), &message)

	}

	// Wait untill wg is 0
	wg.Wait()

	defer os.Exit(0)
}

func connSetup(port string, message *[][]string) {

	conn := listen(port)

	if auto {
		err := InitConn(conn)
		if err != nil {
			log.Println(err)
			return
		}
	}

	utils.Print("Setup finished\n", 1)

	go hostLoop(conn, port, message)

}
