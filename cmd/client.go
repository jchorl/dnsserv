package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jchorl/watchdog"
	"github.com/spf13/cobra"

	"github.com/jchorl/dnsserv/common"
)

var clientCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a DNS entry.",
	Long:  `Update a DNS entry on the server to point to the IP where the request originates from.`,
	Run:   client,
}

var (
	dnsServer string
	domain    string
)

func init() {
	clientCmd.Flags().StringVar(&caPath, "ca-path", "./certs/root.pem", "Path to root ca cert")
	clientCmd.Flags().StringVar(&certPath, "cert-path", "./certs/cert.pem", "Path to cert for mTLS")
	clientCmd.Flags().StringVar(&keyPath, "key-path", "./certs/cert.key", "Path to key for mTLS")
	clientCmd.Flags().StringVar(&dnsServer, "dns-server", "dns.joshchorlton.com", "URL of the DNS server")
	clientCmd.Flags().StringVar(&domain, "domain", "somedomain.com", "Domain to point at this machine")
	rootCmd.AddCommand(clientCmd)
}

func client(cmd *cobra.Command, args []string) {
	tlsConfig := common.LoadTLSConfigOrPanic(caPath, certPath, keyPath)

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	updateRequest := common.UpdateRequest{
		Domain: domain,
	}
	marshaled, err := json.Marshal(updateRequest)
	if err != nil {
		log.Fatalf("Unable to marshal json: %s\n", err)
	}

	resp, err := client.Post(dnsServer, "application/json", bytes.NewReader(marshaled))
	if err != nil {
		log.Fatalf("Error making request: %s\n", err)
	}

	if resp.StatusCode > 204 {
		log.Printf("Received status: %d\n", resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading body: %s\n", err)
		}

		log.Fatalf("Body: %s\n", body)
	}

	wdClient := watchdog.Client{"https://watchdog.joshchorlton.com"}
	wdClient.Ping("dnsserv-pi", watchdog.Watch_DAILY)

	log.Println("Updated successfully")
}
