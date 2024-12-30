package client

import (
	"bufio"
	"log"
	"net"
	"os"
	"time"

	"github.com/adzsx/gwire/internal/crypt"
	"github.com/adzsx/gwire/internal/utils"
)

var (
	conn        net.Conn
	err         error
	inpEnc      string
	displayTime bool
)

// Function for ongoing connection
func client(conn net.Conn) {
	log.SetFlags(0)
	// Receive Data

	utils.Print("Started client routine", 3)

	// Starting ui
	// utils.Ansi("\033[999B")

	// Receive Data
	var data string
	go func() {
		log.SetFlags(0)
		for {
			//Read data
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			data = string(buffer[:bytes])

			if err != nil {
				if err.Error() == "EOF" {
					utils.Print("Connection closed by remote host", 0)
					os.Exit(0)
				}
				log.Fatalln("Error reading data:", err.Error())

			}

			if len([]byte(inpEnc)) != 0 {
				data = crypt.DecryptAES(data, []byte(inpEnc))
			}

			log.Println(data)

		}
	}()

	// Send data
	func() {
		reader := bufio.NewReader(os.Stdin)

		log.SetFlags(0)
		if displayTime {
			log.SetFlags(log.Ltime)
		}

		for {
			time.Sleep(time.Millisecond * 100)

			inp, _ := reader.ReadString('\n')

			//netutils.AddMsg(text)

			if len(inp) > 16384 {
				log.Println("Message cant be over 16384 characters long")
				break
			}

			if len([]byte(inpEnc)) != 0 {
				conn.Write([]byte(crypt.EncryptAES(inp, []byte(inpEnc))))
			} else {
				conn.Write([]byte(inp))
			}

		}
	}()

}
