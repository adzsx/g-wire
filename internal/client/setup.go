package client

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/adzsx/gwire/internal/crypt"
	"github.com/adzsx/gwire/internal/utils"
)

// Function connects to host with TCP
func ClientSetup(host string, port string, enc bool, time bool) {
	log.SetFlags(0)
	if time {
		log.SetFlags(log.Ltime)
	}

	displayTime = time

	// Connect to host

	conn, err = net.Dial("tcp", host+":"+port)

	if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
		utils.Err(errors.New("connection refused by destination"), true)
	}

	utils.Print("Connected to "+host+":"+port+"\n", 0)

	if enc && 1 == 2 {
		err = initClient(conn)
		utils.Err(err, true)
	}

	utils.Print("Setup finished", 1)
	client(conn)

}

// Func for setting up RSA encryption for the clientcs
func initClient(conn net.Conn) error {

	utils.Print("Generating RSA Keys", 1)
	var privateKey rsa.PrivateKey = crypt.GenKeys()

	byteKey := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	utils.Print("Sending Public Key", 2)
	conn.Write(byteKey)

	// Wait for host to send password
	utils.Print("Waiting for response", 2)

	buffer := make([]byte, 512)
	bytes, err := conn.Read(buffer)
	utils.Err(err, true)
	data := buffer[:bytes]

	hexString := hex.EncodeToString(data)

	log.Println(hexString)

	passwd := crypt.DecryptRSA(data, &privateKey)
	utils.Err(err, true)
	inpEnc = passwd

	utils.Print("Password received with length "+strconv.Itoa(len(passwd)), 1)
	utils.Print("Seinding hello world package", 2)

	conn.Write([]byte(crypt.EncryptAES("Hello World", []byte(inpEnc))))

	utils.Print("Sent Control package", 2)

	return nil
}
