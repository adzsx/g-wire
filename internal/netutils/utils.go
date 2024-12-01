package netutils

import (
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/adzsx/gwire/internal/utils"
)

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
