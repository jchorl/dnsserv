package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnsserv",
	Short: "DNSServ is a dynamic dns server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("You must either run with update or serve subcommands")
	},
}

var (
	// flag vars shared by both client and server
	caPath   string
	certPath string
	keyPath  string
)

// Execute is the entrypoint into the app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
