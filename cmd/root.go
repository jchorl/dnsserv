package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "dnsserv",
	Short: "DNSServ is a dynamic dns server",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// For cobra + glog flags. Available to all subcommands.
		goflag.Parse()
	},
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

func init() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

// Execute is the entrypoint into the app
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
