package cmd

import (
	"log"

	"github.com/adzsx/gwire/internal/netutils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "print info about the network",
	Long:  "This command is used to display information about the network currently connected to.",

	Run: info,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().BoolP("public", "p", false, "Request and show public IP from")
	infoCmd.Flags().BoolP("scan", "s", false, "Scan network for active hosts")
}

func info(cmd *cobra.Command, args []string) {
	getPublic, _ := cmd.Flags().GetBool("public")
	scan, _ := cmd.Flags().GetBool("scan")

	log.SetFlags(0)
	privateIP, mask, nHosts := netutils.NetworkInfo()

	if getPublic {
		publicIP := netutils.GetPublicIP()
		log.Printf("Public IP:		%v", publicIP)
	}

	log.Printf("Private IP: 		%v\nSubnetmask: 		%v\nPossible hosts: 	%v", privateIP, mask, nHosts)

	if scan {
		log.Printf("Scanning network...")
	}

}
