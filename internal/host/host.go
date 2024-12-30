package host

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"

	"github.com/adzsx/gwire/internal/crypt"
	"github.com/adzsx/gwire/internal/utils"
)

var (
	enc         string
	displayTime bool
)

// Listener loop for individual port
func hostLoop(conn net.Conn, port string) {

	utils.Print("Started host routine", 3)

	// Read data
	go func() {
		log.SetFlags(0)
		for {
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer

			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			data := string(buffer[:bytes])

			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("Connection on port %v closed", port)
					return
				} else {
					log.Fatalln("Error reading data:", err.Error())
				}

			}

			if len([]byte(enc)) != 0 {
				data = crypt.DecryptAES(data, []byte(enc)) + "\n"
			}

			log.Print(data)

		}

	}()

	// Send data from input
	func() {

		reader := bufio.NewReader(os.Stdin)
		log.SetFlags(0)
		if displayTime {
			log.SetFlags(log.Ltime)
		}

		for {
			//utils.Ansi("\033[999B\033[999D")

			// attach username

			inp, _ := reader.ReadString('\n')

			if inp[len(inp)-1:] == "\n" {
				inp = strings.Replace(inp, "\n", "", 1)
			}

			if len([]byte(enc)) != 0 {
				conn.Write([]byte(crypt.EncryptAES(inp, []byte(enc))))
			} else {
				conn.Write([]byte(inp))
			}
		}
	}()
}
