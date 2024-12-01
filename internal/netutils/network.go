package netutils

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/adzsx/gwire/internal/utils"
)

var (
	colorList []string = []string{"\033[92m", "\033[94m", "\033[95m", "\033[96m"}
	sconn     net.Conn
	counter   int
	found     bool
	accept    bool
)

func Subnet() string {
	cidr, _ := net.InterfaceAddrs()

	return fmt.Sprint(cidr[1])
}

func CalcAddr(cidr string) (string, string) {
	mask := utils.FilterChar(cidr, "/", false)
	maskN, err := strconv.Atoi(mask)

	utils.Err(err, true)

	var subnetmask []string

	for i := 0; i < int(math.Floor(float64(maskN)/8)); i++ {
		subnetmask = append(subnetmask, "255")
	}

	maskN -= len(subnetmask) * 8

	if maskN != 0 {
		converted := 0
		for i := 0; i <= maskN; i++ {
			converted += 256 >> maskN
		}
		log.Println(converted)
		subnetmask = append(subnetmask, strconv.Itoa(converted))
	}

	rest := 4 - len(subnetmask)
	for i := 0; i < rest; i++ {
		subnetmask = append(subnetmask, "0")
	}

	return utils.FilterChar(cidr, "/", true), strings.Join(subnetmask, ".")
}

func GetHosts(cidr string) ([]string, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ipList []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {

		if ip.String()[len(ip.String())-3:] != "255" && ip.String()[len(ip.String())-2:] != ".0" {
			ipList = append(ipList, ip.String())
		}
	}

	return ipList, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func Info() (string, string, string, string) {
	ip, mask := CalcAddr(Subnet())
	list, err := GetHosts(Subnet())
	if err != nil {
		utils.Err(err, true)
	}
	nHosts := strconv.Itoa(len(list))

	return ip, mask, nHosts, ""
}

func AddMsg(msg string, gui bool) {
	fmt.Println()
	utils.Ansi("\033[2A\033[999D\033[K\033[L")

	color := utils.GetRandomString(colorList, utils.FilterChar(msg, ">", true))
	utils.Ansi(color)
	fmt.Print(msg)
	utils.Ansi("\033[999B\033[999D\033[1C")

}

func ScanRange(ips []string, port string) (string, net.Conn) {
	connChan := make(chan net.Conn)

	for _, element := range ips {
		address := element + ":" + port
		counter++
		go scan(address, connChan)
	}

	for counter > 0 {
		time.Sleep(time.Millisecond * 100)
		if counter == 0 {
			break
		}
	}

	if len(connChan) == 0 && !accept {
		utils.Err(errors.New("no host found"), true)
		os.Exit(0)
	}

	sconn = <-connChan

	ip := utils.FilterChar(sconn.RemoteAddr().String(), ":", true)

	return ip, sconn
}

func scan(address string, connChan chan net.Conn) {
	ping := Ping(utils.FilterChar(address, ":", true))
	if !ping {
		time.Sleep(time.Millisecond)
		counter--
		return
	}
	conn, err := net.Dial("tcp", address)

	if err != nil {
		/* log.Println(counter) */
		time.Sleep(time.Millisecond)
		counter--
		return
	}

	for found {
		if !found {
			break
		}
		time.Sleep(time.Second)
	}

	for {
		if !found {
			found = true
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Found open port on %v\nDo you want to connect? [y/n] ", utils.FilterChar(address, ":", true))
			input, _ := reader.ReadString('\n')

			input = input[0:1]

			if input == "n" || input == "no" {

				counter--
				if counter == 0 {
					return
				}
				reader = bufio.NewReader(os.Stdin)
				fmt.Print("Continue scan anyways? [y/n] ")
				input, _ := reader.ReadString('\n')

				input = input[0:1]

				log.Println()

				found = false

				if input == "n" || input == "no" {
					conn.Close()
					os.Exit(0)
				} else {
					return

				}
			}
			accept = true

			counter = 0

			connChan <- conn
		}

		time.Sleep(time.Second)

	}

}

func Ping(ip string) bool {

	cmd := exec.Command("ping", "-i", "0.2", "-c", "3", "-w", "1", ip)
	out, _ := cmd.Output()

	output, _ := strconv.Atoi(utils.FilterChar(utils.FilterChar(string(out), ",", false), ",", true)[1:2])

	return output > 1

}
