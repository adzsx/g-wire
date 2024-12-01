package netutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/adzsx/gwire/internal/utils"
)

type IP struct {
	Query string
}

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

func NetworkInfo() (string, string, string) {
	ip, mask := CalcAddr(Subnet())
	list, err := GetHosts(Subnet())
	if err != nil {
		utils.Err(err, true)
	}
	nHosts := strconv.Itoa(len(list))

	return ip, mask, nHosts
}

func AddMsg(msg string, gui bool) {
	fmt.Println()
	utils.Ansi("\033[2A\033[999D\033[K\033[L")

	color := utils.GetRandomString(colorList, utils.FilterChar(msg, ">", true))
	utils.Ansi(color)
	fmt.Print(msg)
	utils.Ansi("\033[999B\033[999D\033[1C")

}

func Ping(ip string) bool {

	cmd := exec.Command("ping", "-i", "0.2", "-c", "3", "-w", "1", ip)
	out, _ := cmd.Output()

	output, _ := strconv.Atoi(utils.FilterChar(utils.FilterChar(string(out), ",", false), ",", true)[1:2])

	return output > 1

}

func GetPublicIP() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}
