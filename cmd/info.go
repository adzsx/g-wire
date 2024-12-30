package cmd

import (
	"log"

	"github.com/adzsx/gwire/internal/netutils"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "get info about the network",
	Long:  "This command is used to display information about the network currently connected to.",
	Run:   scan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	scanCmd.Flags().BoolP("scan", "s", false, "Scan network for active hosts")
}

func scan(cmd *cobra.Command, args []string) {
	scan, _ := cmd.Flags().GetBool("scan")

	log.SetFlags(0)
	privateIP, mask, nHosts := netutils.NetworkInfo()

	log.Printf("Private IP: 		%v\nSubnetmask: 		%v\nPossible hosts: 	%v", privateIP, mask, nHosts)

	if scan {
		log.Printf("Scanning network...")
	}

}
