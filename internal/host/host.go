package host

import (
	"bufio"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/adzsx/gwire/internal/netutils"
	"github.com/adzsx/gwire/internal/utils"
	"github.com/adzsx/gwire/pkg/crypt"
)

var (
	auto bool
	wg   sync.WaitGroup
	sent int

	ports []string
	//handshake   bool
	enc         string
	username    string
	displayTime bool
)

// Set up connection to specific port
func listen(port string) net.Conn {

	log.SetFlags(0)
	if displayTime {
		log.SetFlags(log.Ltime)
	}

	//Listen and connect
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		log.Fatalln("Permission denied.\nTry again with root or take a port above 1023")
		wg.Done()
		os.Exit(0)
	}

	utils.Print("Listening on port "+port, 1)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalln("Error accepting connection:", err.Error())
		wg.Done()
		return conn
	}

	utils.Print("Connected to "+conn.RemoteAddr().String(), 0)

	return conn
}

// Listener loop for individual port
func hostLoop(conn net.Conn, port string, message *[][]string) {

	utils.Print("Started host routine", 3)

	fmt.Println()

	// Read data
	var receivedHost []string
	go func() {
		for {
			//Make buffer for read data
			buffer := make([]byte, 16384)
			//Write length of message to bytes, message to buffer
			bytes, err := conn.Read(buffer)
			// Iterate for length over message
			data := string(buffer[:bytes])
			receivedHost = append(receivedHost, data)

			if err != nil {
				if err.Error() == "EOF" {
					log.Printf("Connection on port %v closed", port)
					wg.Done()
					return
				} else {
					log.Fatalln("Error reading data:", err.Error())
				}

			}

			if len(port) > 1 {
				*message = append(*message, []string{utils.FilterChar(conn.LocalAddr().String(), ":", false), data})
			}
		}

	}()

	// Function for printing messages
	go func() {
		var data string
		for {
			//time.Sleep(time.Millisecond * time.Duration(input.TimeOut))

			if len(receivedHost) != 0 {

				if len([]byte(enc)) != 0 {
					data = crypt.DecryptAES(receivedHost[0], []byte(enc)) + "\n"
				} else {
					data = receivedHost[0] + "\n"
				}

				netutils.AddMsg(data, false)

				receivedHost = utils.Remove(receivedHost, receivedHost[0])
			}
		}
	}()

	// Send data from input
	go func() {

		reader := bufio.NewReader(os.Stdin)
		log.SetFlags(0)
		if displayTime {
			log.SetFlags(log.Ltime)
		}

		for {
			utils.Ansi("\033[999B\033[999D")
			fmt.Print(">")

			// attach username
			text := username + "> "
			inp, _ := reader.ReadString('\n')

			text += inp

			utils.Ansi("\033[1A\033[K")

			netutils.AddMsg(text, false)

			if text[len(text)-1:] == "\n" {
				text = strings.Replace(text, "\n", "", 1)
			}

			if len(port) > 1 {

				if len(enc) != 0 {

					*message = append(*message, []string{"0", crypt.EncryptAES(text, []byte(enc))})
				} else {
					*message = append(*message, []string{"0", text})
				}

				sent = -1

			} else {
				if len([]byte(enc)) != 0 {
					conn.Write([]byte(crypt.EncryptAES(text, []byte(enc))))
				} else {
					conn.Write([]byte(text))
				}
			}
		}
	}()

	//send data from other clients
	if len(port) > 1 {
		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				if len(*message) > 0 {

					for _, element := range *message {
						if element[0] != utils.FilterChar(conn.LocalAddr().String(), ":", false) {
							conn.Write([]byte(element[1]))
							sent += 1
						}
						time.Sleep(time.Millisecond * 50)
					}
					if sent == len(port)-1 {
						*message = [][]string{}
						sent = 0
					}

				}
			}
		}()
	}
}

func InitConn(conn net.Conn) error {
	// Make buffer for receiving RSA public key
	utils.Print("Waiting for RSA key from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)
	buffer := make([]byte, 4096)
	bytes, err := conn.Read(buffer)
	if err != nil {
		wg.Done()
		return errors.New("connection closed")
	}
	sentPublicKey := buffer[:bytes]

	// Convert bytes back to public key
	publicKey, err := x509.ParsePKCS1PublicKey(sentPublicKey)

	if err != nil {
		wg.Done()
		return errors.New("received data not RSA publickey")
	}

	utils.Print("Publickey received from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)

	// Send encrypted AES key over connection
	utils.Print("Sending Password", 2)
	encKey := crypt.EncryptRSA(*publicKey, []byte(enc))
	conn.Write(encKey)

	utils.Print("Waiting for password confirmation", 2)

	buffer = make([]byte, 512)
	bytes, err = conn.Read(buffer)
	utils.Err(err, true)
	data := string(buffer[:bytes])

	utils.Print("Received password confirmation", 2)

	if crypt.DecryptAES(data, []byte(enc)) != enc {
		conn.Write([]byte("wrong password"))
		return errors.New("wrong password")
	}

	conn.Write([]byte(crypt.EncryptAES("success", []byte(enc))))

	return nil
}
