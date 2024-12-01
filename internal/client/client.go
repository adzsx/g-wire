package client

import (
	"bufio"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/adzsx/gwire/internal/netutils"
	"github.com/adzsx/gwire/internal/utils"
	"github.com/adzsx/gwire/pkg/crypt"
)

var (
	conn net.Conn
	err  error

	inpHost      string
	inpPort      string
	inpHandshake bool
	inpEnc       string
	inpUsername  string
	displayTime  bool
)

// Function for ongoing connection
func client(conn net.Conn) {
	// Receive Data

	utils.Print("Started client routine", 3)

	// Starting ui
	utils.Ansi("\033[999B")

	// Receive Data
	var received []string

	go func() {
		for {
			//Read data
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			received = append(received, string(buffer[:bytes]))

			if err != nil {
				if err.Error() == "EOF" {
					utils.Print("Connection closed by remote host", 0)
					os.Exit(0)
				}
				log.Fatalln("Error reading data:", err.Error())

			}

		}
	}()

	// Function for printing received data
	go func() {
		var data string

		for {
			//time.Sleep(time.Millisecond * time.Duration(input.TimeOut))

			if len(received) != 0 {
				if len([]byte(inpEnc)) != 0 {
					data = crypt.DecryptAES(received[0], []byte(inpEnc))
				} else {
					data = received[0]
				}

				netutils.AddMsg(data, false)

				received = utils.Remove(received, received[0])

			}

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
			fmt.Print("\033[999B\033[999D")
			fmt.Print(">")
			time.Sleep(time.Millisecond * 100)

			text := inpUsername + "> "

			inp, _ := reader.ReadString('\n')
			text += inp

			utils.Ansi("\033[1A\033[K")

			netutils.AddMsg(text, false)

			if len(text) > 16384 {
				log.Println("Message cant be over 16384 characters long")
				break
			}

			if len([]byte(inpEnc)) != 0 {
				conn.Write([]byte(crypt.EncryptAES(text, []byte(inpEnc))))
			} else {
				conn.Write([]byte(text))
			}

		}
	}()

}

// Func for setting up RSA encryption for the clientcs
func initClient(conn net.Conn) error {

	utils.Print("Generating RSA Keys", 1)
	var rsaKeys = crypt.GenKeys()

	byteKey := x509.MarshalPKCS1PublicKey(&rsaKeys.PublicKey)

	utils.Print("Sending Public Key", 2)
	conn.Write(byteKey)

	// Wait for host to send password
	utils.Print("Waiting for response", 2)
	buffer := make([]byte, 512)
	bytes, err := conn.Read(buffer)
	if err != nil {
		log.Println("Connection unexpectedly closed. Aborting Setup")
		return errors.New("connection failed")
	}
	data := string(buffer[:bytes])

	passwd := crypt.DecryptRSA(rsaKeys, []byte(data))

	inpEnc = passwd

	utils.Print("Password received", 1)
	utils.Print("Seinding Password confirmation package", 2)

	conn.Write([]byte(crypt.EncryptAES(inpEnc, []byte(inpEnc))))

	utils.Print("Waiting for Control package", 2)

	buffer = make([]byte, 512)
	bytes, err = conn.Read(buffer)
	utils.Err(err, true)
	data = string(buffer[:bytes])
	data = crypt.DecryptAES(data, []byte(inpEnc))

	utils.Print("Received Control Package", 2)

	if string(data) == "wrong password" {
		log.Println("Wrong password. Aborting connection")
		return errors.New("wrong password")
	}

	return nil
}
