package host

import (
	"crypto/x509"
	"errors"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/adzsx/gwire/internal/crypt"
	"github.com/adzsx/gwire/internal/utils"
)

var (
	conn net.Conn
)

// Set up listener for each port on list
func HostSetup(port int, src string, enc bool, time bool) {
	log.SetFlags(0)
	displayTime = time
	// Global slice for distributing messages

	sPort := strconv.Itoa(port)

	//Listen and connect
	ln, err := net.Listen("tcp", ":"+sPort)

	if err != nil && strings.Contains(err.Error(), "permission denied") {
		log.Fatalln("Permission denied.\nTry again with root or take a port above 1023")
		os.Exit(0)
	}

	utils.Print("Listening on port "+sPort, 1)

	for {

		conn, err = ln.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection:", err.Error())
		}

		if src != strings.Split(conn.RemoteAddr().String(), ":")[0] && src != "" {
			utils.Print("Rejecting "+conn.RemoteAddr().String()+" (Not matching src)", 2)
			conn.Close()
			continue
		}

		break
	}

	utils.Print("Connected to "+conn.RemoteAddr().String(), 0)

	if enc && 1 == 2 {
		err := InitConn(conn)
		if err != nil {
			log.Println(err)
			return
		}
	}

	utils.Print("Setup finished\n", 1)

	hostLoop(conn, sPort)
}

func InitConn(conn net.Conn) error {
	// Make buffer for receiving RSA public key
	utils.Print("Waiting for RSA key from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)
	buffer := make([]byte, 4096)
	bytes, err := conn.Read(buffer)
	if err != nil {
		return errors.New("connection closed")
	}
	sentPublicKey := buffer[:bytes]

	// Convert bytes back to public key
	publicKey, err := x509.ParsePKCS1PublicKey(sentPublicKey)

	if err != nil {
		return errors.New("received data not RSA publickey")
	}

	utils.Print("Publickey received from "+utils.FilterChar(conn.RemoteAddr().String(), ":", true)+"\n", 1)

	// Send encrypted AES key over connection
	utils.Print("Sending Password", 2)

	encKey := crypt.EncryptRSA([]byte(enc), publicKey)
	conn.Write(encKey)
	utils.Print("Sent password. Package length: "+strconv.Itoa(len(encKey)), 2)

	utils.Print("Waiting for control package", 2)

	buffer = make([]byte, 88)
	_, err = conn.Read(buffer)
	utils.Err(err, true)
	data := string(buffer)

	utils.Print("Received control package", 2)

	if crypt.DecryptAES(data, []byte(enc)) != enc {
		conn.Write([]byte("wrong password"))
		return errors.New("wrong password")
	}

	conn.Write([]byte(crypt.EncryptAES("success", []byte(enc))))

	return nil
}
